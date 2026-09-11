package workflow

import (
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"picup/core/domain/constants"
	"picup/core/domain/models"
	"picup/core/domain/repositories"
	"picup/core/domain/schema"
	"picup/core/infrastructure/config"
	"picup/core/infrastructure/exifworker"
	"picup/core/usecases/variant"
	"strings"

	"github.com/google/uuid"
)

type Workflow struct {
	Files       FileStorage
	ExifDataMap exifworker.ExifDataMap
	Repo        repositories.Repository

	imMaster         *models.Image
	imMasterFilepath string

	images    []*models.Image
	filepaths []string
}

func (c *Workflow) GetMasterCopy() (*models.Image, error) {
	if c.imMaster == nil {
		return nil, fmt.Errorf("master image is not initialized")
	}

	// 1. Shallow copy the top-level struct
	cp := *c.imMaster
	return &cp, nil
}

// ProcessMaster is the main entry point.
func (c *Workflow) ProcessMaster(ctx context.Context, file io.Reader, originalFilename string) error {
	fwtask, err := c.Files.Stage(file, filepath.Ext(originalFilename))
	if err != nil {
		return fmt.Errorf("failed to stage file to temp: %w", err)
	}
	defer fwtask.Clean()

	tmpFilepath := fwtask.WorkFilepath()

	exifdata, err := c.inspectMetadata(tmpFilepath)
	if err != nil {
		return fmt.Errorf("failed to inspect metadata: %w", err)
	}

	dv := config.Variant.GetMaster()
	im := models.NewImage(constants.IMAGE_TYPE_DEFAULT, dv.Name, originalFilename)

	im.MIMEType = exifdata.MimeType()
	im.Width = exifdata.ImageWidth()
	im.Height = exifdata.ImageHeight()
	im.SizeBytes = exifdata.FileSizeBytes()
	im.StoredName = fmt.Sprintf("%s.%s", uuid.New().String(), exifdata.FileExt())

	finalPath, err := c.Files.Save(fwtask, filepath.Join(im.Variant, im.StoredName))
	if err != nil {
		return fmt.Errorf("save master file: %w", err)
	}

	c.imMasterFilepath = finalPath

	c.images = append(c.images, im)
	c.filepaths = append(c.filepaths, c.imMasterFilepath)

	c.imMaster = im
	return nil
}

func (c *Workflow) ProcessDisplay(ctx context.Context, dv *schema.Variant) error {
	if c.imMaster == nil {
		return fmt.Errorf("cannot process display variant: master image is not initialized")
	}

	fwtask, err := c.Files.Copy(c.imMasterFilepath)
	if err != nil {
		return fmt.Errorf("prepare display work file: %w", err)
	}
	defer fwtask.Clean()

	workftype := filepath.Ext(c.imMasterFilepath)
	newWorkftype := fmt.Sprintf(".%s", dv.Type)
	workfpath := fwtask.WorkFilepath()
	newWorkfpath := fwtask.NewWorkFilepath(fmt.Sprintf("%s%s", uuid.New().String(), newWorkftype))

	if err := variant.Process(ctx, variant.ProcessRequest{
		VariantName:  dv.Name,
		WorkFilePath: workfpath,
		WorkFileType: workftype,
		DestFilePath: newWorkfpath,
		DestFileType: newWorkftype,
		SourceWidth:  c.imMaster.Width,
		SourceHeight: c.imMaster.Height,
	}); err != nil {
		return fmt.Errorf("failed to process display variant %q: %w", dv.Name, err)
	}

	exifdata, err := c.inspectMetadata(fwtask.WorkFilepath())
	if err != nil {
		return fmt.Errorf("failed to inspect display variant metadata: %w", err)
	}

	im := models.NewImage(constants.IMAGE_TYPE_DISPLAY, dv.Name, c.imMaster.StoredName)

	im.MasterID = &c.imMaster.ID
	im.MIMEType = c.imMaster.MIMEType
	im.Width = c.imMaster.Width
	im.Height = c.imMaster.Height
	im.SizeBytes = c.imMaster.SizeBytes
	im.StoredName = fmt.Sprintf("%s.%s", uuid.New().String(), exifdata.FileExt())

	targetFilepath, err := c.Files.Save(fwtask, filepath.Join(im.Variant, im.StoredName))
	if err != nil {
		return fmt.Errorf("save display file: %w", err)
	}

	c.images = append(c.images, im)
	c.filepaths = append(c.filepaths, targetFilepath)

	return nil
}

func (c *Workflow) ProcessVariant(ctx context.Context, v *schema.Variant) error {
	if v.IsDisplay() || v.IsMaster() {
		return nil
	}

	if c.imMaster == nil {
		return fmt.Errorf("cannot process variant %q: master image is not initialized", v.Name)
	}

	fwtask, err := c.Files.Copy(c.imMasterFilepath)
	if err != nil {
		return fmt.Errorf("prepare variant work file: %w", err)
	}
	defer fwtask.Clean()

	workftype := filepath.Ext(c.imMasterFilepath)
	newWorkftype := fmt.Sprintf(".%s", v.Type)
	workfpath := fwtask.WorkFilepath()
	newWorkfpath := fwtask.NewWorkFilepath(fmt.Sprintf("%s%s", uuid.New().String(), newWorkftype))

	if err := variant.Process(ctx, variant.ProcessRequest{
		VariantName:  v.Name,
		WorkFilePath: workfpath,
		WorkFileType: workftype,
		DestFilePath: newWorkfpath,
		DestFileType: newWorkftype,
		SourceWidth:  c.imMaster.Width,
		SourceHeight: c.imMaster.Height,
	}); err != nil {
		return fmt.Errorf("failed to process variant %q: %w", v.Name, err)
	}

	exifdata, err := c.inspectMetadata(fwtask.WorkFilepath())
	if err != nil {
		return fmt.Errorf("failed to inspect variant %q metadata: %w", v.Name, err)
	}

	im := models.NewImage(constants.IMAGE_TYPE_DEFAULT, v.Name, c.imMaster.StoredName)

	im.MasterID = &c.imMaster.ID
	im.MIMEType = exifdata.MimeType()
	im.Width = exifdata.ImageWidth()
	im.Height = exifdata.ImageHeight()
	im.SizeBytes = exifdata.FileSizeBytes()
	im.StoredName = fmt.Sprintf("%s.%s", uuid.New().String(), exifdata.FileExt())

	targetFilepath, err := c.Files.Save(fwtask, filepath.Join(im.Variant, im.StoredName))
	if err != nil {
		return fmt.Errorf("save variant file: %w", err)
	}

	c.images = append(c.images, im)
	c.filepaths = append(c.filepaths, targetFilepath)

	return nil
}

func (c *Workflow) Finalize(ctx context.Context) error {
	// NOTE: not optimized yet
	for _, im := range c.images {
		log.Printf("saving image %s variant %s", im.ID, im.Variant)
		if err := c.Repo.SaveImage(ctx, im); err != nil {
			return fmt.Errorf("failed to save variant %q of image %q to repository: %w", im.Variant, im.ID, err)
		}
	}
	return nil
}

func (c *Workflow) Rollback(ctx context.Context) {
	if c.imMaster != nil {
		_ = c.Repo.DeleteImage(ctx, c.imMaster.ID)
	}
	if c.Files != nil {
		_ = c.Files.Delete(c.filepaths...)
	}
}

func (c *Workflow) inspectMetadata(sourceFilepath string) (*exifworker.Metadata, error) {
	sourceFilename := filepath.Base(sourceFilepath)

	if err := exifworker.ReadFile(c.ExifDataMap, sourceFilepath); err != nil {
		return nil, fmt.Errorf("failed to read exif data from %q: %w", sourceFilename, err)
	}

	exifdata, exists := c.ExifDataMap[sourceFilename]
	if !exists || exifdata == nil {
		return nil, fmt.Errorf("metadata entry missing for file %q", sourceFilename)
	}

	if !strings.HasPrefix(exifdata.MimeType(), "image") {
		return nil, fmt.Errorf("unsupported media file type %q for file %q", exifdata.MimeType(), sourceFilename)
	}

	return exifdata, nil
}
