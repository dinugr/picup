package exifworker

import (
	"log"
	"sync/atomic"
)

// Global package variable using atomic.Pointer for safe concurrent access
var defaultInstance atomic.Pointer[ExifWorker]

// SetDefault updates the default global worker instance.
func SetDefault(instance *ExifWorker, err error) *ExifWorker {
	if err != nil {
		log.Fatalf("Failed to initialize exif worker: %v", err)
	}
	defaultInstance.Store(instance)
	return instance
}

func ReadFile(result ExifDataMap, filePath string) error {
	return defaultInstance.Load().ReadFile(result, filePath)
}

func Instance() *ExifWorker {
	return defaultInstance.Load()
}
