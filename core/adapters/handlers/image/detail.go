package imagehandler

import (
	"fmt"
	"net/http"
	"picup/core/adapters/handlers/dto"
	"picup/core/domain/constants"
	"picup/core/domain/models"
	"strings"
)

func (h *ImageHandler) GetDetailByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if strings.TrimSpace(id) == "" {
		h.writeError(w, http.StatusBadRequest, "Missing image id")
		return
	}

	imlist, _, err := h.repo.ListImage(r.Context(), models.ImageQuery{
		ID:        id,
		Type:      []string{constants.IMAGE_TYPE_DEFAULT},
		PageIndex: constants.PAGE_INDEX_DISABLE,
	})

	if err != nil {
		h.writeError(w, http.StatusNotFound, fmt.Sprintf("Image not found: %s", id))
		return
	}

	master := dto.Image{}
	variants := []dto.Image{}

	for _, v := range imlist {
		imdto := h.toImageDetailResponse(r, v)
		variants = append(variants, imdto)

		if v.IsMaster() {
			master = imdto
		}
	}

	h.writeJSON(w, http.StatusOK, dto.ImageDetailResponse{
		Success: true,
		Data: dto.ImageDetail{
			Master:   master,
			Variants: variants,
			Notes:    "", // STUB: assign notes after adding feature!
		},
	})
}
