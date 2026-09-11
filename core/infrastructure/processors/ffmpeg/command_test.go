package ffmpeg

import (
	"context"
	"reflect"
	"testing"

	"picup/core/domain/variantargs"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		name       string
		spec       variantargs.Pipeline
		sourceType string
		want       []string
	}{
		{
			name: "translate includes quality and preserves operation order",
			spec: variantargs.Pipeline{Operations: []variantargs.Operation{
				variantargs.Scale{Size: dimension(256, 256), Quality: 72},
				variantargs.Crop{Size: dimension(256, 256), Gravity: centerGravity()},
			}},
			sourceType: "jpg",
			want:       []string{"-vf", "scale=256:256:force_original_aspect_ratio=decrease,crop=256:256:(iw-out_w)/2:(ih-out_h)/2", "-quality", "72"},
		},
		{
			name: "translate for heic",
			spec: variantargs.Pipeline{Operations: []variantargs.Operation{
				variantargs.Scale{Size: dimension(256, 256), Quality: 72},
			}},
			sourceType: "heic",
			want:       []string{"-filter_complex", "[0:g:0]scale=256:256:force_original_aspect_ratio=decrease[out]", "-map", "[out]", "-quality", "72"},
		},
		{
			name: "translate for gif",
			spec: variantargs.Pipeline{Operations: []variantargs.Operation{
				variantargs.Scale{Size: dimension(900, 900), Quality: 80, Mode: variantargs.ScaleModeFixed},
			}},
			sourceType: "gif",
			want:       []string{"-vf", "scale=900:900", "-quality", "80", "-loop", "0"},
		},
		{
			name:       "translate uses default extension handling",
			spec:       variantargs.Pipeline{Operations: []variantargs.Operation{variantargs.Scale{Size: dimension(128, 0), Quality: 80}}},
			sourceType: "raw",
			want:       []string{"-vf", "scale=128:-1", "-quality", "80"},
		},
		{
			name:       "empty operations",
			spec:       variantargs.Pipeline{},
			sourceType: "jpg",
			want:       []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command, err := translatePipelineToCommand(testContext{pipeline: test.spec, workType: test.sourceType})
			if err != nil {
				t.Fatalf("pipelineToCommand() error = %v", err)
			}
			if !reflect.DeepEqual(command, test.want) {
				t.Fatalf("pipelineToCommand() = %#v, want %#v", command, test.want)
			}
		})
	}
}

func TestTranslateCropGravity(t *testing.T) {
	tests := []struct {
		name    string
		gravity variantargs.CropDirection
		want    string
	}{
		{name: "north west", gravity: variantargs.CropDirection{Vertical: variantargs.DirectionNorth, Horizontal: variantargs.DirectionWest}, want: "crop=128:64:0:0"},
		{name: "south east", gravity: variantargs.CropDirection{Vertical: variantargs.DirectionSouth, Horizontal: variantargs.DirectionEast}, want: "crop=128:64:iw-out_w:ih-out_h"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := translateCrop(variantargs.Crop{Size: dimension(128, 64), Gravity: test.gravity})
			if err != nil {
				t.Fatalf("translateCrop() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("translateCrop() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestTranslateRejectsInvalidGravity(t *testing.T) {
	_, err := translatePipelineToCommand(testContext{pipeline: variantargs.Pipeline{Operations: []variantargs.Operation{
		variantargs.Crop{Size: dimension(128, 128), Gravity: variantargs.CropDirection{Vertical: "x", Horizontal: variantargs.DirectionCenter}},
	}}, workType: "jpg"})
	if err == nil {
		t.Fatal("pipelineToCommand() error = nil, want invalid gravity error")
	}
}

func TestTranslateScaleDimensionsAndModes(t *testing.T) {
	tests := []struct {
		name string
		size variantargs.Dimension
		mode variantargs.ScaleMode
		want string
	}{
		{name: "width only", size: dimension(128, 0), want: "scale=128:-1"},
		{name: "height only", size: dimension(0, 64), want: "scale=-1:64"},
		{name: "shrink", size: dimension(128, 64), mode: variantargs.ScaleModeLockRatioShrink, want: "scale=128:64:force_original_aspect_ratio=decrease"},
		{name: "extend", size: dimension(128, 64), mode: variantargs.ScaleModeLockRatioExtend, want: "scale=128:64:force_original_aspect_ratio=increase"},
		{name: "fixed", size: dimension(128, 64), mode: variantargs.ScaleModeFixed, want: "scale=128:64"},
		{name: "ignored dimension", size: variantargs.Dimension{}, want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := translateScale(variantargs.Scale{Size: test.size, Mode: test.mode})
			if err != nil {
				t.Fatalf("translateScale() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("translateScale() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestTranslateCropPartialDimensions(t *testing.T) {
	tests := []struct {
		name string
		size variantargs.Dimension
		want string
	}{
		{name: "width only", size: dimension(128, 0), want: "crop=128:ih:(iw-out_w)/2:(ih-out_h)/2"},
		{name: "height only", size: dimension(0, 64), want: "crop=iw:64:(iw-out_w)/2:(ih-out_h)/2"},
		{name: "ignored dimension", size: variantargs.Dimension{}, want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := translateCrop(variantargs.Crop{Size: test.size, Gravity: centerGravity()})
			if err != nil {
				t.Fatalf("translateCrop() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("translateCrop() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestTranslateRejectsUnsupportedOperation(t *testing.T) {
	_, err := translatePipelineToCommand(testContext{
		pipeline: variantargs.Pipeline{Operations: []variantargs.Operation{nil}},
		workType: "jpg",
	})
	if err == nil {
		t.Fatal("pipelineToCommand() error = nil, want unsupported operation error")
	}
}

func TestProcessReportsPipelineTranslationError(t *testing.T) {
	ctx := testContext{
		pipeline: variantargs.Pipeline{Operations: []variantargs.Operation{
			variantargs.Crop{
				Size: dimension(128, 128),
				Gravity: variantargs.CropDirection{
					Vertical:   "invalid",
					Horizontal: variantargs.DirectionCenter,
				},
			},
		}},
		workType: "jpg",
	}

	err := (&FFMPEGProcessor{}).Process(ctx)
	if err == nil {
		t.Fatal("Process() error = nil, want pipeline translation error")
	}
	want := `translator failed: translate variant arguments: unsupported vertical crop gravity "invalid"`
	if err.Error() != want {
		t.Fatalf("Process() error = %q, want %q", err, want)
	}
}

type testContext struct {
	pipeline variantargs.Pipeline
	workType string
}

func (c testContext) Context() context.Context       { return context.Background() }
func (c testContext) Pipeline() variantargs.Pipeline { return c.pipeline }
func (c testContext) WorkFileType() string           { return c.workType }
func (testContext) DestFileType() string             { return "jpg" }
func (testContext) WorkFilePath() string             { return "source.jpg" }
func (testContext) DestFilePath() string             { return "destination.jpg" }
func (testContext) SourceWidth() int                 { return 0 }
func (testContext) SourceHeight() int                { return 0 }

func dimension(width, height int) variantargs.Dimension {
	var result variantargs.Dimension
	if width > 0 {
		result.Width = &width
	}
	if height > 0 {
		result.Height = &height
	}
	return result
}

func centerGravity() variantargs.CropDirection {
	return variantargs.CropDirection{Vertical: variantargs.DirectionCenter, Horizontal: variantargs.DirectionCenter}
}
