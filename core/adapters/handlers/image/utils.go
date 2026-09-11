package imagehandler

import (
	"encoding/json"
	"net/http"
	"path"
	"picup/core/adapters/handlers/dto"
	"picup/core/domain/models"
	"picup/core/infrastructure/utils"
)

func (h *ImageHandler) toImageDetailResponse(r *http.Request, img *models.Image) dto.Image {

	return dto.Image{
		ID:         img.ID,
		Name:       img.UploadName,
		StoredName: img.StoredName,
		Variant:    img.Variant,
		MIMEType:   img.MIMEType,
		SizeBytes:  img.SizeBytes,
		Width:      img.Width,
		Height:     img.Height,
		CreatedAt:  img.CreatedAt,
		ImageURL:   utils.ResolveImageURL(path.Join(img.Variant, img.StoredName)),
	}
}

func (h *ImageHandler) toImageListItemResponse(r *http.Request, img *models.Image) dto.ImageItem {

	return dto.ImageItem{
		ID:        *img.MasterID,
		Name:      img.UploadName,
		SizeBytes: img.SizeBytes,
		CreatedAt: img.CreatedAt,
		ImageUrl:  utils.ResolveImageURL(path.Join(img.Variant, img.StoredName)),
	}
}

func (h *ImageHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *ImageHandler) writeError(w http.ResponseWriter, status int, msg string) {
	h.writeJSON(w, status, dto.ErrorResponse{
		Success: false,
		Error:   msg,
	})
}
