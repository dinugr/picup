package imagehandler

import (
	"net/http"

	"picup/core/domain/repositories"
	"picup/core/infrastructure/exifworker"
	"picup/core/usecases/workflow"
)

type ImageHandler struct {
	repo        repositories.Repository
	exifworker  *exifworker.ExifWorker
	fileStorage workflow.FileStorage
}

func NewImageHandler(
	repo repositories.Repository,
	fileStorage workflow.FileStorage,
) *ImageHandler {
	return &ImageHandler{
		repo:        repo,
		fileStorage: fileStorage,
	}
}

func (h *ImageHandler) RouteInit(mux *http.ServeMux, chain func(http.Handler) http.Handler) {
	mux.Handle("GET /api/images", chain(http.HandlerFunc(h.ListImages)))
	mux.Handle("POST /api/images", chain(http.HandlerFunc(h.UploadImage)))
	mux.Handle("GET /api/images/{id}", chain(http.HandlerFunc(h.GetDetailByID)))
	mux.Handle("DELETE /api/images/{id}", chain(http.HandlerFunc(h.DeleteImage)))
}
