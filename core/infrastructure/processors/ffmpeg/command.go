package ffmpeg

import (
	"fmt"
	"strconv"
	"strings"

	"picup/core/domain/variantargs"
	"picup/core/usecases/variant"
)

func translatePipelineToCommand(context variant.ProcessContext) ([]string, error) {
	filters := make([]string, 0, len(context.Pipeline().Operations))
	quality := 0

	for _, operation := range context.Pipeline().Operations {
		switch value := operation.(type) {
		case variantargs.Scale:
			filter, err := translateScale(value)
			if err != nil {
				return []string{}, fmt.Errorf("translate variant arguments: %w", err)
			}
			if filter != "" {
				filters = append(filters, filter)
			}
			quality = value.Quality
		case variantargs.Crop:
			filter, err := translateCrop(value)
			if err != nil {
				return []string{}, fmt.Errorf("translate variant arguments: %w", err)
			}
			if filter != "" {
				filters = append(filters, filter)
			}
		default:
			return []string{}, fmt.Errorf("unsupported variant operation %T", operation)
		}
	}

	if len(filters) == 0 {
		return []string{}, nil
	}

	filterGraph := strings.Join(filters, ",")
	sourceType := context.WorkFileType()
	args := make([]string, 0, 6)

	switch sourceType {
	case "heic":
		args = append(args, "-filter_complex", "[0:g:0]"+filterGraph+"[out]", "-map", "[out]")
	default:
		args = append(args, "-vf", filterGraph)
	}

	if quality > 0 {
		args = append(args, "-quality", strconv.Itoa(quality))
	}
	if sourceType == "gif" {
		args = append(args, "-loop", "0")
	}

	return args, nil
}

func normalizeExtension(value string) string {
	return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(value), "."))
}

func translateScale(value variantargs.Scale) (string, error) {
	width, height := translateDimensions(value.Size, "-1")
	if width == "-1" && height == "-1" {
		return "", nil
	}

	filter := "scale=" + width + ":" + height
	if value.Size.Width == nil || value.Size.Height == nil {
		return filter, nil
	}

	switch value.Mode {
	case variantargs.ScaleModeFixed:
		return filter, nil
	case variantargs.ScaleModeLockRatioExtend:
		return filter + ":force_original_aspect_ratio=increase", nil
	case "", variantargs.ScaleModeLockRatioShrink:
		return filter + ":force_original_aspect_ratio=decrease", nil
	default:
		return "", fmt.Errorf("unsupported scale mode %q", value.Mode)
	}
}

func translateCrop(value variantargs.Crop) (string, error) {
	width := "iw"
	height := "ih"
	if value.Size.Width != nil {
		width = strconv.Itoa(*value.Size.Width)
	}
	if value.Size.Height != nil {
		height = strconv.Itoa(*value.Size.Height)
	}
	if width == "iw" && height == "ih" {
		return "", nil
	}

	x := "(iw-out_w)/2"
	y := "(ih-out_h)/2"
	switch value.Gravity.Horizontal {
	case variantargs.DirectionWest:
		x = "0"
	case variantargs.DirectionEast:
		x = "iw-out_w"
	case variantargs.DirectionCenter:
	default:
		return "", fmt.Errorf("unsupported horizontal crop gravity %q", value.Gravity.Horizontal)
	}
	switch value.Gravity.Vertical {
	case variantargs.DirectionNorth:
		y = "0"
	case variantargs.DirectionSouth:
		y = "ih-out_h"
	case variantargs.DirectionCenter:
	default:
		return "", fmt.Errorf("unsupported vertical crop gravity %q", value.Gravity.Vertical)
	}

	return "crop=" + width + ":" + height + ":" + x + ":" + y, nil
}

func translateDimensions(dimension variantargs.Dimension, empty string) (string, string) {
	width := empty
	height := empty
	if dimension.Width != nil {
		width = strconv.Itoa(*dimension.Width)
	}
	if dimension.Height != nil {
		height = strconv.Itoa(*dimension.Height)
	}
	return width, height
}
