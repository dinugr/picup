package dto

import "time"

type Image struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	StoredName string    `json:"stored_name"`
	Variant    string    `json:"variant"`
	MIMEType   string    `json:"mime_type"`
	SizeBytes  int64     `json:"size_bytes"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	CreatedAt  time.Time `json:"created_at"`
	ImageURL   string    `json:"image_url"`
}

type ImageItem struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	SizeBytes int64     `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
	ImageUrl  string    `json:"image_url"`
}

type ImageDetail struct {
	Master   Image   `json:"master"`
	Variants []Image `json:"variants"`
	Notes    string  `json:"notes"`
}

type ImageResponse struct {
	Success bool  `json:"success"`
	Data    Image `json:"data"`
}

type ImageListDefaultResponse struct {
	Success    bool           `json:"success"`
	Data       []ImageItem    `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

type ImageListVariantResponse struct {
	Success    bool           `json:"success"`
	Data       []Image        `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

type ImageDetailResponse struct {
	Success bool        `json:"success"`
	Data    ImageDetail `json:"data"`
}
