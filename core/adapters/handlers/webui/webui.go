package webuihandler

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"picup/core/infrastructure/config"
	"strings"
)

//go:embed default.html
var defaultHTML string // Unexported; kept out of public package API

type WebUIHandler struct {
}

func NewWebUIHandler() *WebUIHandler {
	return &WebUIHandler{}
}

func (h *WebUIHandler) RouteInit(mux *http.ServeMux, chain func(http.Handler) http.Handler) {
	frontendDist := h.resolveCandidate()

	if frontendDist != "" {
		log.Printf("Serving static frontend files from %s", frontendDist)
		fileServer := http.FileServer(http.Dir(frontendDist))
		mux.Handle("/", chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Join(frontendDist, r.URL.Path)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				http.ServeFile(w, r, filepath.Join(frontendDist, "index.html"))
				return
			}
			fileServer.ServeHTTP(w, r)
		})))
	} else {
		replacer := strings.NewReplacer(
			"{{port}}", fmt.Sprintf("%d", config.Server.Port),
			"{{maxupload}}", fmt.Sprintf("%d", config.Server.MaxSize/(1024*1024)),
		)

		mux.Handle("/", chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, replacer.Replace(defaultHTML))
		})))
	}
}

func (h *WebUIHandler) resolveCandidate() string {
	if config.Server.WebUIEnabled {
		return config.Server.WebUIDir
	}
	return ""
}
