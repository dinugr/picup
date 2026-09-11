package variant

import (
	"context"
	"picup/core/domain/schema"
	"picup/core/domain/variantargs"
)

type ProcessContext interface {
	Context() context.Context
	Pipeline() variantargs.Pipeline

	WorkFileType() string // ext
	DestFileType() string // ext
	WorkFilePath() string // path/to/temp/filename
	DestFilePath() string // path/to/temp/filename

	SourceWidth() int
	SourceHeight() int
}

type ProcessRequest struct {
	VariantName string

	WorkFilePath string
	WorkFileType string
	DestFilePath string
	DestFileType string

	SourceWidth  int
	SourceHeight int
}

type VariantProcessor interface {
	Supports(variant *schema.Variant) error
	Process(context ProcessContext) error
}
