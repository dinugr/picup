package config

import (
	"os"
	"path/filepath"
	"testing"

	"picup/core/domain/schema"
)

func TestCompileVariants(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantOps int
		wantErr bool
	}{
		{name: "valid arguments", input: "scale=256:72,crop=256", wantOps: 2},
		{name: "empty arguments", input: "", wantOps: 0},
		{name: "invalid arguments", input: "scale=256:101", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			variants := schema.VariantMap{
				"thumbnail": {
					Name:    "thumbnail",
					RawArgs: test.input,
				},
			}

			err := CompileVariants(variants)
			if test.wantErr && err == nil {
				t.Fatal("CompileVariants() error = nil")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("CompileVariants() error = %v", err)
			}
			if !test.wantErr && len(variants["thumbnail"].Pipeline.Operations) != test.wantOps {
				t.Fatalf("unexpected operation count: got %d, want %d", len(variants["thumbnail"].Pipeline.Operations), test.wantOps)
			}
		})
	}
}

func TestCompileVariantsAggregatesErrors(t *testing.T) {
	variants := schema.VariantMap{
		"bad-scale": {RawArgs: "scale=256:101"},
		"bad-crop":  {RawArgs: "crop=256:nn"},
	}

	if err := CompileVariants(variants); err == nil {
		t.Fatal("CompileVariants() error = nil")
	}
}

func TestLoadExternalToolPaths(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.ini")
	configContents := `[server]
file_types = png
exiftool_path = /opt/bin/exiftool

[engine.custom]
name = custom
provider = ffmpeg
path = /opt/bin/ffmpeg

[variant.preview]
type = webp
engine = custom
`
	if err := os.WriteFile(configPath, []byte(configContents), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	Server = schema.Server{}
	Engine = schema.EngineMap{}
	Variant = schema.VariantMap{}

	if err := Load(configPath); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got := Server.ExifToolPath; got != "/opt/bin/exiftool" {
		t.Fatalf("Server.ExifToolPath = %q, want %q", got, "/opt/bin/exiftool")
	}
	if got := Engine["custom"].Parameters["path"]; got != "/opt/bin/ffmpeg" {
		t.Fatalf("engine path = %q, want %q", got, "/opt/bin/ffmpeg")
	}
}
