package variantargs

import (
	"fmt"
	"strconv"
	"strings"
)

type ValueField struct {
	value string
}

func NewValueField(value string) ValueField {
	return ValueField{value: value}
}

func (field ValueField) AsStr() (string, error) {
	return field.value, nil
}

func (field ValueField) AsStrDefault(defaultValue string) (string, error) {
	if field.value == "" {
		return defaultValue, nil
	}
	return field.AsStr()
}

func (field ValueField) AsInt() (int, error) {
	return strconv.Atoi(field.value)
}

func (field ValueField) AsIntDefault(defaultValue int) (int, error) {
	if field.value == "" {
		return defaultValue, nil
	}
	return field.AsInt()
}

func (field ValueField) AsFloat() (float64, error) {
	return strconv.ParseFloat(field.value, 64)
}

func (field ValueField) AsFloatDefault(defaultValue float64) (float64, error) {
	if field.value == "" {
		return defaultValue, nil
	}
	return field.AsFloat()
}

func (field ValueField) AsDimension() (Dimension, error) {
	value := field.value
	if value == "" || value == "x" {
		return Dimension{}, nil
	}

	if !strings.Contains(value, "x") {
		size, err := parsePositiveInteger(value)
		if err != nil {
			return Dimension{}, fmt.Errorf("expected a positive integer or <width>x<height>")
		}
		return Dimension{Width: &size, Height: &size}, nil
	}

	parts := strings.Split(value, "x")
	if len(parts) != 2 || (parts[0] == "" && parts[1] == "") {
		return Dimension{}, fmt.Errorf("expected <width>x<height> with at least one dimension")
	}
	var dimension Dimension
	if parts[0] != "" {
		size, err := parsePositiveInteger(parts[0])
		if err != nil {
			return Dimension{}, fmt.Errorf("width must be a positive integer")
		}
		dimension.Width = &size
	}
	if parts[1] != "" {
		size, err := parsePositiveInteger(parts[1])
		if err != nil {
			return Dimension{}, fmt.Errorf("height must be a positive integer")
		}
		dimension.Height = &size
	}
	return dimension, nil
}

func (field ValueField) AsDimensionDefault(defaultValue Dimension) (Dimension, error) {
	if field.value == "" {
		return defaultValue, nil
	}
	return field.AsDimension()
}

func parsePositiveInteger(value string) (int, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return 0, fmt.Errorf("not positive")
	}
	return parsed, nil
}
