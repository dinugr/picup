package variantargs

import (
	"reflect"
	"testing"
)

func TestValueFieldDefaults(t *testing.T) {
	tests := []struct {
		name          string
		field         ValueField
		wantString    string
		wantInt       int
		wantFloat     float64
		wantDimension Dimension
	}{
		{
			name:          "empty field",
			field:         NewValueField(""),
			wantString:    "fallback",
			wantInt:       7,
			wantFloat:     1.5,
			wantDimension: Dimension{Width: new(10), Height: new(20)},
		},
		{
			name:          "non-empty field",
			field:         NewValueField("42"),
			wantString:    "42",
			wantInt:       42,
			wantFloat:     42,
			wantDimension: Dimension{Width: new(42), Height: new(42)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotString, stringErr := test.field.AsStrDefault("fallback")
			if stringErr != nil || gotString != test.wantString {
				t.Fatalf("AsStrDefault() = %q, %v; want %q, nil", gotString, stringErr, test.wantString)
			}

			gotInt, intErr := test.field.AsIntDefault(7)
			if intErr != nil || gotInt != test.wantInt {
				t.Fatalf("AsIntDefault() = %d, %v; want %d, nil", gotInt, intErr, test.wantInt)
			}

			gotFloat, floatErr := test.field.AsFloatDefault(1.5)
			if floatErr != nil || gotFloat != test.wantFloat {
				t.Fatalf("AsFloatDefault() = %v, %v; want %v, nil", gotFloat, floatErr, test.wantFloat)
			}

			gotDimension, dimensionErr := test.field.AsDimensionDefault(test.wantDimension)
			if dimensionErr != nil || !reflect.DeepEqual(gotDimension, test.wantDimension) {
				t.Fatalf("AsDimensionDefault() = %#v, %v; want %#v, nil", gotDimension, dimensionErr, test.wantDimension)
			}
		})
	}
}

func TestValueFieldDefaultsPreserveConversionErrors(t *testing.T) {
	field := NewValueField("invalid")

	if _, err := field.AsIntDefault(7); err == nil {
		t.Fatal("AsIntDefault() error = nil, want conversion error")
	}
	if _, err := field.AsFloatDefault(1.5); err == nil {
		t.Fatal("AsFloatDefault() error = nil, want conversion error")
	}
	if _, err := field.AsDimensionDefault(Dimension{}); err == nil {
		t.Fatal("AsDimensionDefault() error = nil, want conversion error")
	}
}
