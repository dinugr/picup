package ffmpeg

import (
	"bytes"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"picup/core/domain/schema"
	"picup/core/usecases/variant"
)

type FFMPEGProcessor struct {
	ffmpegExecutablePath string
}

func init() {
	variant.RegisterFactory("ffmpeg", New)
}

func New(fctx *variant.FactoryContext) (variant.VariantProcessor, error) {
	executablePath := "ffmpeg"
	if fctx != nil && fctx.Parameters != nil && fctx.Parameters["path"] != "" {
		executablePath = fctx.Parameters["path"]
	}

	resolvedPath, err := exec.LookPath(executablePath)
	if err != nil {
		return nil, fmt.Errorf("ffmpeg executable %q is not installed or not in PATH: %w", executablePath, err)
	}

	return &FFMPEGProcessor{ffmpegExecutablePath: resolvedPath}, nil
}

func (j *FFMPEGProcessor) Supports(variant *schema.Variant) error {
	supportedResultType := "png|jpg|jpeg|jxl|webp|avif|gif|ico|heic|bmp|tiff"
	targetType := strings.ToLower(variant.Type)

	if slices.Contains(strings.Split(supportedResultType, "|"), targetType) {
		return nil
	}

	return fmt.Errorf("unsupported variant type: %s", variant.Type)
}

func (j *FFMPEGProcessor) Process(ctx variant.ProcessContext) error {
	commands, err := translatePipelineToCommand(ctx)
	if err != nil {
		return fmt.Errorf("translator failed: %w", err)
	}

	args := []string{}
	args = append(args, "-y")                     // dont ask, just approve !
	args = append(args, "-loglevel", "warning")   // only show warning log
	args = append(args, "-i", ctx.WorkFilePath()) // source
	args = append(args, commands...)              // doing filter stuff
	args = append(args, "-map_metadata", "-1")    // strip metadata
	args = append(args, ctx.DestFilePath())       // result

	// Run the FFmpeg command.
	cmd := exec.CommandContext(ctx.Context(), j.ffmpegExecutablePath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w, error details: %s", err, stderr.String())
	}

	return nil
}
