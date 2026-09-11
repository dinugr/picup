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
	baseurl        string
	datadir        string
	assetsBasePath string
}

var defaultInstance atomic.Pointer[UrlUtils]

func NewBaseURL(host string, port uint64, datadir string, baseurl string) string {

	addr := fmt.Sprintf("%s:%d", host, port)
	if host == "*" || host == "" || host == "0.0.0.0" {
		addr = fmt.Sprintf(":%d", port)
	}

	if baseurl == "" {
		baseurl = fmt.Sprintf("http://%s/", addr)
	}

	// Normalize baseurl to never end with a trailing slash.
	// This makes URL joining predictable.
	baseurl = strings.TrimSuffix(baseurl, "/")

	instance := &UrlUtils{
		baseurl:        baseurl,
		datadir:        filepath.Clean(datadir),
		assetsBasePath: strings.TrimSuffix(strings.TrimSpace(config.Server.AssetsBasePath), "/"),
	}

	defaultInstance.Store(instance)

	return addr

}

func ResolveImageURL(sourceFilepath string) string {
	instance := defaultInstance.Load()
	p, _ := strings.CutPrefix(filepath.Clean(sourceFilepath), instance.datadir)

	// instance.baseurl is normalized to never end with '/'.
	// Join relative path segments without introducing double slashes.
	// Note: AssetsBasePath is expected to be a URL path prefix (e.g. "/assets/images").
	return instance.baseurl + path.Join(instance.assetsBasePath, p)
}
