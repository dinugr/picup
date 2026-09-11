package dto

type UploadConfig struct {
	SizeLimitBytes   int64    `json:"size_limit_bytes"`
	AllowedFileTypes []string `json:"allowed_file_types"`
}

type VariantConfig struct {
	Name      string `json:"name"`
	Format    string `json:"format"`
	Path      string `json:"path"`
	Arguments string `json:"arguments"`
}

type ConfigResponse struct {
	Upload           UploadConfig             `json:"upload"`
	Target           map[string]VariantConfig `json:"target"`
	VariantsSequence []string                 `json:"variants_sequence"`
}
