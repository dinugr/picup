package imagehandler

import (
	"fmt"
	"log"
	"net/http"
	"picup/core/adapters/handlers/dto"
	"picup/core/infrastructure/config"
	"picup/core/infrastructure/exifworker"
	"picup/core/usecases/workflow"
)

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {

	maxSize := config.Server.MaxSize

	if err := r.ParseMultipartForm(maxSize + (1024 * 1024)); err != nil {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse multipart form: %v", err))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Missing or invalid 'file' in multipart form data")
		return
	}
	defer file.Close()

	exifDataMap := exifworker.ExifDataMap{}

	wf := &workflow.Workflow{
		ExifDataMap: exifDataMap,
		Files:       h.fileStorage,
		Repo:        h.repo,
	}

	if err := wf.ProcessMaster(r.Context(), file, header.Filename); err != nil {
		wf.Rollback(r.Context())
		log.Printf("processing stream error : %v", err)
		h.writeError(w, http.StatusInternalServerError, "Unexpected error")
		return
	}

	if err := wf.ProcessDisplay(r.Context(), config.Variant.GetDisplay()); err != nil {
		wf.Rollback(r.Context())
		log.Printf("processing display error : %v", err)
		h.writeError(w, http.StatusInternalServerError, "Unexpected error")
		return
	}

	for _, v := range config.Variant {
		if err := wf.ProcessVariant(r.Context(), v); err != nil {
			log.Printf("[WARN] failed to process variant %s: %v", v.Key, err.Error())
		}
	}

	if err := wf.Finalize(r.Context()); err != nil {
		wf.Rollback(r.Context())
		log.Printf("processing display error : %v", err)
		h.writeError(w, http.StatusInternalServerError, "Unexpected error")
		return
	}

	master, _ := wf.GetMasterCopy()
	itemResp := h.toImageDetailResponse(r, master)
	h.writeJSON(w, http.StatusCreated, dto.ImageResponse{
		Success: true,
		Data:    itemResp,
	})

}
