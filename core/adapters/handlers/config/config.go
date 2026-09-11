package confighandler

import (
	"encoding/json"
	"net/http"
	"strings"

	"picup/core/adapters/handlers/dto"
	"picup/core/infrastructure/config"
)

type ConfigHandler struct {
}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

func (h *ConfigHandler) RouteInit(mux *http.ServeMux, chain func(http.Handler) http.Handler) {
	mux.Handle("GET /api/config", chain(http.HandlerFunc(h.GetConfig)))
}

func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sizeLimitBytes := config.Server.MaxSize
	allowedTypesStr := config.Server.FileTypes

	targetPresets := make(map[string]dto.VariantConfig)
	for name, item := range config.Variant {
		targetPresets[name] = dto.VariantConfig{
			Name:   item.Name,
			Format: item.Type,
			Path:   item.Name,
			//Arguments: item.Arguments,
		}
	}

	var allowedTypes []string
	for _, t := range allowedTypesStr {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			allowedTypes = append(allowedTypes, trimmed)
		}
	}

	res := dto.ConfigResponse{
		Upload: dto.UploadConfig{
			SizeLimitBytes:   int64(sizeLimitBytes),
			AllowedFileTypes: allowedTypes,
		},
		Target:           targetPresets,
		VariantsSequence: config.Server.VariantsSequence,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
