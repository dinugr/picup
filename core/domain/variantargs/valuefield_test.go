package variantargs

import (
	"reflect"
	"testing"
)

func TestValueFieldAsStr(t *testing.T) {
	field := NewValueField("value")
	got, err := field.AsStr()
	if err != nil {
		t.Fatalf("AsStr() error = %v", err)
	}
	if got != "value" {
		t.Fatalf("AsStr() = %q, want value", got)
	}
}

func TestValueFieldAsInt(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    int
		wantErr bool
	}{
		{name: "positive", value: "42", want: 42},
		{name: "negative", value: "-7", want: -7},
		{name: "invalid", value: "nope", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewValueField(test.value).AsInt()
			if (err != nil) != test.wantErr {
				t.Fatalf("AsInt() error = %v, wantErr %v", err, test.wantErr)
			}
			if !test.wantErr && got != test.want {
				t.Fatalf("AsInt() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestValueFieldAsFloat(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    float64
		wantErr bool
	}{
		{name: "decimal", value: "1.25", want: 1.25},
		{name: "integer", value: "7", want: 7},
		{name: "invalid", value: "nope", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewValueField(test.value).AsFloat()
			if (err != nil) != test.wantErr {
				t.Fatalf("AsFloat() error = %v, wantErr %v", err, test.wantErr)
			}
			if !test.wantErr && got != test.want {
				t.Fatalf("AsFloat() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestValueFieldAsDimension(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    Dimension
		wantErr bool
	}{
		{name: "empty", value: "", want: Dimension{}},
		{name: "empty dimensions", value: "x", want: Dimension{}},
		{name: "square", value: "128", want: Dimension{Width: new(128), Height: new(128)}},
		{name: "width only", value: "128x", want: Dimension{Width: new(128)}},
		{name: "height only", value: "x64", want: Dimension{Height: new(64)}},
		{name: "rectangle", value: "128x64", want: Dimension{Width: new(128), Height: new(64)}},
		{name: "zero", value: "0", wantErr: true},
		{name: "malformed", value: "12xx34", wantErr: true},
		{name: "missing dimensions", value: "xx", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewValueField(test.value).AsDimension()
			if (err != nil) != test.wantErr {
				t.Fatalf("AsDimension() error = %v, wantErr %v", err, test.wantErr)
			}
			if !test.wantErr && !reflect.DeepEqual(got, test.want) {
				t.Fatalf("AsDimension() = %#v, want %#v", got, test.want)
			}
		})
	}
}
