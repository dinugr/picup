package schema

type Server struct {
	Host string
	Port uint64

	// AssetsBasePath is the URL path prefix under which image variants are served.
	// Example: "/assets/images" (no trailing slash).
	AssetsBasePath string

	DataDir      string
	TempDir      string
	MaxSize      int64
	FileTypes    []string
	ExifToolPath string

	VariantsSequence []string

	// webui
	WebUIEnabled bool
	WebUIDir     string

	// auth
	Auth Auth

	// database
	DSN string // driver_name://arguments...
}
