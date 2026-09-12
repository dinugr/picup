package utils

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"

	"picup/core/infrastructure/config"
)

type UrlUtils struct {
	datadir        string
	assetsBasePath string
}

var defaultInstance atomic.Pointer[UrlUtils]

func NewBaseURL(host string, port uint64, datadir string) string {
	addr := fmt.Sprintf("%s:%d", host, port)
	if host == "*" || host == "" || host == "0.0.0.0" {
		addr = fmt.Sprintf(":%d", port)
	}

	instance := &UrlUtils{
		datadir:        filepath.Clean(datadir),
		assetsBasePath: strings.TrimSuffix(strings.TrimSpace(config.Server.AssetsBasePath), "/"),
	}

	defaultInstance.Store(instance)

	return addr
}

func ResolveImageURL(sourceFilepath string) string {
	instance := defaultInstance.Load()
	p, _ := strings.CutPrefix(filepath.Clean(sourceFilepath), instance.datadir)

	// AssetsBasePath is expected to be a URL path prefix (e.g. "/assets/images").
	return path.Join(instance.assetsBasePath, p)
}
