package imagehandler

import (
	"fmt"
	"net/http"
	"picup/core/adapters/handlers/dto"
	"picup/core/domain/constants"
	"picup/core/domain/models"
	"strconv"
)

func (h *ImageHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	search := query.Get("search")

	images, totalItems, err := h.repo.ListImage(r.Context(), models.ImageQuery{
		Type: []string{constants.IMAGE_TYPE_DISPLAY},

		Search:    search,
		PageIndex: page,
		PageSize:  limit,
		SortBy:    "created_at",
		SortOrder: constants.SORT_ORDER_DESC,
	})
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to query images: %v", err))
		return
	}

	totalPages := (totalItems + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	var items []dto.ImageItem
	for _, img := range images {
		items = append(items, h.toImageListItemResponse(r, img))
	}

	if items == nil {
		items = []dto.ImageItem{}
	}

	h.writeJSON(w, http.StatusOK, dto.ImageListDefaultResponse{
		Success: true,
		Data:    items,
		Pagination: dto.PaginationInfo{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	})
}
