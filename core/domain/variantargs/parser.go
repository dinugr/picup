package variantargs

import (
	"fmt"
	"strings"
)

const (
	defaultQuality       = 80
	defaultInterpolation = "bicubic"
)

const (
	delimiterGroup = ","
	delimiterField = "="
	delimiterValue = ":"
)

// <fieldname1>=<args1>[;<fieldname2>=<args1>,<args2>,<args3>[;...]]
func Parse(input string) (Pipeline, error) {
	if strings.TrimSpace(input) == "" {
		return Pipeline{}, nil
	}

	parts := strings.Split(input, delimiterGroup)
	operations := make([]Operation, 0, len(parts))
	errors := make([]error, 0)
	seen := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		name, value, found := strings.Cut(part, delimiterField)
		name = strings.TrimSpace(name)
		if !found || name == "" {
			errors = append(errors, &ParseError{Field: "operation", Value: part, Reason: "expected name=value"})
			continue
		}
		if _, exists := seen[name]; exists {
			errors = append(errors, &ParseError{Operation: name, Field: "operation", Value: name, Reason: "duplicate operation"})
			continue
		}
		rawValueFields := strings.Split(value, delimiterValue)
		valuefields := make([]ValueField, len(rawValueFields))
		for index, rawValueField := range rawValueFields {
			valuefields[index] = NewValueField(rawValueField)
		}

		seen[name] = struct{}{}

		var operation Operation
		var err error

		switch name {
		case "scale":
			operation, err = parseScale(valuefields...)
		case "crop":
			operation, err = parseCrop(valuefields...)
		default:
			err = &ParseError{Operation: name, Field: "operation", Value: name, Reason: "unsupported operation"}
		}
		if err != nil {
			errors = append(errors, err)
			continue
		}
		operations = append(operations, operation)
	}

	if len(errors) > 0 {
		return Pipeline{Operations: operations}, &ParseErrors{Errors: errors}
	}

	return Pipeline{Operations: operations}, nil
}

// <size>:<quality>:<mode>:<interpolation>
func parseScale(fields ...ValueField) (Scale, error) {
	if len(fields) > 4 {
		return Scale{}, fieldError("scale", "value", "", "expected at most four fields")
	}
	fields = append(fields, NewValueField(""), NewValueField(""), NewValueField(""))

	size, err := fields[0].AsDimension()
	if err != nil {
		return Scale{}, operationError("scale", "size", fields[0].value, err)
	}

	quality, err := fields[1].AsIntDefault(defaultQuality)
	if err != nil || quality < 1 || quality > 100 {
		return Scale{}, fieldError("scale", "quality", fields[1].value, "must be an integer between 1 and 100")
	}

	mode := ScaleModeLockRatioShrink
	if fields[2].value != "" {
		mode = ScaleMode(fields[2].value)
		switch mode {
		case ScaleModeFixed, ScaleModeLockRatioExtend, ScaleModeLockRatioShrink:
		default:
			return Scale{}, fieldError("scale", "mode", fields[2].value, "must be fixed, lrext, or lrshr")
		}
	}

	interpolation, err := fields[3].AsStrDefault(defaultInterpolation)
	if err != nil {
		return Scale{}, operationError("scale", "interpolation", fields[3].value, err)
	}

	return Scale{Size: size, Quality: quality, Mode: mode, Interpolation: interpolation}, nil
}

// <size>:<gravity>
func parseCrop(fields ...ValueField) (Crop, error) {
	if len(fields) > 2 {
		return Crop{}, fieldError("crop", "value", "", "expected at most two fields")
	}
	fields = append(fields, NewValueField(""))

	size, err := fields[0].AsDimension()
	if err != nil {
		return Crop{}, operationError("crop", "size", fields[0].value, err)
	}

	gravity := CropDirection{Vertical: DirectionCenter, Horizontal: DirectionCenter}
	if fields[1].value != "" {
		gravityValue, conversionErr := fields[1].AsStr()
		if conversionErr != nil {
			return Crop{}, operationError("crop", "gravity", fields[1].value, conversionErr)
		}
		gravity, err = parseGravity(gravityValue)
		if err != nil {
			return Crop{}, operationError("crop", "gravity", fields[1].value, err)
		}
	}

	return Crop{Size: size, Gravity: gravity}, nil
}

func parseGravity(value string) (CropDirection, error) {
	if value == "c" {
		return CropDirection{Vertical: DirectionCenter, Horizontal: DirectionCenter}, nil
	}
	if len(value) != 2 {
		return CropDirection{}, fmt.Errorf("expected c or two-character vertical/horizontal gravity")
	}
	vertical := Direction(value[0])
	horizontal := Direction(value[1])
	if vertical != DirectionNorth && vertical != DirectionSouth && vertical != DirectionCenter {
		return CropDirection{}, fmt.Errorf("invalid vertical direction")
	}
	if horizontal != DirectionWest && horizontal != DirectionEast && horizontal != DirectionCenter {
		return CropDirection{}, fmt.Errorf("invalid horizontal direction")
	}
	return CropDirection{Vertical: vertical, Horizontal: horizontal}, nil
}

func operationError(operation, field, value string, err error) error {
	return fieldError(operation, field, value, err.Error())
}

func fieldError(operation, field, value, reason string) error {
	return &ParseError{Operation: operation, Field: field, Value: value, Reason: reason}
}
