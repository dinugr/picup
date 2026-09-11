package imagehandler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"picup/core/adapters/handlers/dto"
	"picup/core/infrastructure/config"
	"strings"
)

func (h *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.writeError(w, http.StatusBadRequest, "Missing image id")
		return
	}

	imgs, err := h.repo.ListImageByID(r.Context(), id)

	if err != nil || len(imgs) == 0 {
		h.writeError(w, http.StatusNotFound, fmt.Sprintf("Image not found: %s", id))
		return
	}

	if err := h.repo.DeleteImage(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete image metadata: %v", err))
		return
	}

	fnames := []string{}
	for _, item := range imgs {
		fnames = append(fnames, filepath.Join(config.Server.DataDir, item.Variant, item.StoredName))
	}

	// Remove physical file
	if err := h.fileStorage.Delete(fnames...); err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete image files: %v", err))
		return
	}

	h.writeJSON(w, http.StatusOK, dto.MessageResponse{
		Success: true,
		Message: "Image deleted successfully",
	})
}
