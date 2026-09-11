package assets

import (
	"net/http"
	"path"
	"path/filepath"
	"picup/core/infrastructure/config"
)

type AssetsHandler struct {
}

func NewAssetsHandler() *AssetsHandler {
	return &AssetsHandler{}
}

func (h *AssetsHandler) RouteInit(mux *http.ServeMux, chain func(http.Handler) http.Handler) {
	for _, variant := range config.Variant {
		pathURL := path.Join(config.Server.AssetsBasePath, variant.Name) + "/"
		pathDir := filepath.Join(config.Server.DataDir, variant.Name)
		fileServer := http.StripPrefix(pathURL, http.FileServer(http.Dir(pathDir)))

		if variant.IsReserved() {
			mux.Handle(pathURL, chain(fileServer))
			continue
		}
		mux.Handle(pathURL, fileServer)
	}
}
