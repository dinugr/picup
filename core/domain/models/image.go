package models

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UploadName string    `json:"upload_name"` // filename with extension
	StoredName string    `json:"stored_name"` // filename with extension
	Type       string    `json:"type"`        // current possible values: list, detail
	Variant    string    `json:"variant"`
	MasterID   *string   `json:"master_id"`
	MIMEType   string    `json:"mime_type"`
	SizeBytes  int64     `json:"size_bytes"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
}

func (i *Image) IsMaster() bool {
	return i.MasterID == nil || i.MasterID == &i.ID
}

func NewImage(image_type string, variant string, uploadFilename string) *Image {
	return &Image{
		ID:         uuid.New().String(),
		CreatedAt:  time.Now().UTC(),
		Type:       image_type,
		Variant:    variant,
		UploadName: uploadFilename,
	}
}

type ImageQuery struct {
	ID          string   // default ""
	Type        []string // default []
	TypeExclude bool     // default false

	Search    string // default ""
	SortBy    string // default created_at
	SortOrder string // default asc
	PageIndex int    // start from 0
	PageSize  int    // min item = 10
}
