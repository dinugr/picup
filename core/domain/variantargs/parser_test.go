package variantargs

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseScale(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Scale
	}{
		{name: "defaults", input: "scale=128", want: Scale{Size: Dimension{Width: new(128), Height: new(128)}, Quality: 80, Mode: ScaleModeLockRatioShrink, Interpolation: "bicubic"}},
		{name: "fixed mode", input: "scale=128x:100:fixed:linear", want: Scale{Size: Dimension{Width: new(128)}, Quality: 100, Mode: ScaleModeFixed, Interpolation: "linear"}},
		{name: "extend mode", input: "scale=128x:100:lrext:linear", want: Scale{Size: Dimension{Width: new(128)}, Quality: 100, Mode: ScaleModeLockRatioExtend, Interpolation: "linear"}},
		{name: "empty optional fields", input: "scale=128:::bilinear", want: Scale{Size: Dimension{Width: new(128), Height: new(128)}, Quality: 80, Mode: ScaleModeLockRatioShrink, Interpolation: "bilinear"}},
		{name: "shrink mode", input: "scale=128x:100:lrshr:linear", want: Scale{Size: Dimension{Width: new(128)}, Quality: 100, Mode: ScaleModeLockRatioShrink, Interpolation: "linear"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Parse(test.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if !reflect.DeepEqual(got.Operations, []Operation{test.want}) {
				t.Fatalf("Parse() = %#v, want %#v", got.Operations, []Operation{test.want})
			}
		})
	}
}

func TestParseCropAndOrder(t *testing.T) {
	got, err := Parse("scale=128,crop=x64:nw")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	want := []Operation{
		Scale{Size: Dimension{Width: new(128), Height: new(128)}, Quality: 80, Mode: ScaleModeLockRatioShrink, Interpolation: "bicubic"},
		Crop{Size: Dimension{Height: new(64)}, Gravity: CropDirection{Vertical: DirectionNorth, Horizontal: DirectionWest}},
	}
	if !reflect.DeepEqual(got.Operations, want) {
		t.Fatalf("Parse() = %#v, want %#v", got.Operations, want)
	}
}

func TestParseEmptyDimensionsAndCenterGravity(t *testing.T) {
	got, err := Parse("scale=x,crop=:")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if got.Operations[0].(Scale).Size != (Dimension{}) {
		t.Fatalf("scale dimension = %#v, want ignored dimension", got.Operations[0].(Scale).Size)
	}
	crop := got.Operations[1].(Crop)
	if crop.Gravity != (CropDirection{Vertical: DirectionCenter, Horizontal: DirectionCenter}) {
		t.Fatalf("crop gravity = %#v, want center", crop.Gravity)
	}
}

func TestParseRejectsInvalidInput(t *testing.T) {
	tests := []string{
		"scale=0",
		"scale=128:0",
		"scale=128:101",
		"scale=128::2",
		"scale=128::::linear",
		"crop=128:n",
		"crop=128:neast",
		"crop=128:nn",
		"unknown=128",
		"scale=128,scale=256",
		"scale",
		"scale=12xx",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			var parseErr *ParseError
			_, err := Parse(input)
			if err == nil || !errors.As(err, &parseErr) {
				t.Fatalf("Parse() error = %v, want ParseError", err)
			}
		})
	}
}

func TestParseCollectsErrorsAndValidOperations(t *testing.T) {
	got, err := Parse("scale=128,crop=128:nn,unknown=1,scale=256")
	if err == nil {
		t.Fatal("Parse() error = nil, want aggregate errors")
	}

	var parseErrors *ParseErrors
	if !errors.As(err, &parseErrors) {
		t.Fatalf("Parse() error = %T, want ParseErrors", err)
	}
	if len(parseErrors.Errors) != 3 {
		t.Fatalf("error count = %d, want 3", len(parseErrors.Errors))
	}
	if len(got.Operations) != 1 {
		t.Fatalf("operation count = %d, want valid operations preserved", len(got.Operations))
	}
	if _, ok := got.Operations[0].(Scale); !ok {
		t.Fatalf("operation type = %T, want Scale", got.Operations[0])
	}
}

func TestParseEmptyInput(t *testing.T) {
	got, err := Parse("  ")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(got.Operations) != 0 {
		t.Fatalf("operations = %#v, want empty", got.Operations)
	}
}

func TestOperationKinds(t *testing.T) {
	if got := (Scale{}).opKind(); got != "scale" {
		t.Fatalf("Scale.operationKind() = %q, want scale", got)
	}
	if got := (Crop{}).opKind(); got != "crop" {
		t.Fatalf("Crop.operationKind() = %q, want crop", got)
	}
}

func TestParseErrorFormatting(t *testing.T) {
	parseErr := &ParseError{Operation: "scale", Field: "quality", Value: "101", Reason: "out of range"}
	if got, want := parseErr.Error(), `variantargs: operation "scale", field "quality", value "101": out of range`; got != want {
		t.Fatalf("ParseError.Error() = %q, want %q", got, want)
	}

	aggregateErr := &ParseErrors{Errors: []error{parseErr, &ParseError{Operation: "crop", Field: "gravity", Value: "nn", Reason: "invalid"}}}
	if got, want := aggregateErr.Error(), `variantargs: operation "scale", field "quality", value "101": out of range; variantargs: operation "crop", field "gravity", value "nn": invalid`; got != want {
		t.Fatalf("ParseErrors.Error() = %q, want %q", got, want)
	}
	if len(aggregateErr.Unwrap()) != 2 {
		t.Fatalf("ParseErrors.Unwrap() length = %d, want 2", len(aggregateErr.Unwrap()))
	}
}

func TestParseRejectsRemainingMalformedFields(t *testing.T) {
	tests := []string{
		"crop=128:nw:extra",
		"scale=12x34x56",
		"scale=0x12",
		"scale=x0",
		"scale=abc",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := Parse(input)
			if err == nil {
				t.Fatalf("Parse(%q) error = nil, want ParseError", input)
			}
			var parseErr *ParseError
			if !errors.As(err, &parseErr) {
				t.Fatalf("Parse(%q) error = %T, want ParseError", input, err)
			}
		})
	}
}

func TestParseGravityDirections(t *testing.T) {
	valid := []string{"se", "sc", "cw", "ce", "cc", "c"}
	for _, gravity := range valid {
		t.Run("valid_"+gravity, func(t *testing.T) {
			if _, err := Parse("crop=128:" + gravity); err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
		})
	}

	invalid := []string{
		"aw", "sa", // non-constant
		"ec", "en", // flipped
	}
	for _, gravity := range invalid {
		t.Run("invalid_"+gravity, func(t *testing.T) {
			if _, err := Parse("crop=128:" + gravity); err == nil {
				t.Fatalf("Parse() error = nil, want invalid gravity error")
			}
		})
	}
}

func TestParseCropValidationBranches(t *testing.T) {
	if _, err := Parse("crop=invalid:c"); err == nil {
		t.Fatal("Parse() error = nil, want invalid crop size error")
	}
	if _, err := Parse("crop=128:c"); err != nil {
		t.Fatalf("Parse() error = %v, want center gravity to be valid", err)
	}
}
