package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	httpadapter "picup/core/adapters/http"
	"picup/core/infrastructure/config"
	"picup/core/infrastructure/db"
	"picup/core/infrastructure/exifworker"
	"picup/core/infrastructure/utils"
	"picup/core/usecases/variant"
	"picup/core/usecases/workflow"
	"syscall"

	_ "picup/core/infrastructure/processors/ffmpeg"
)

var version = "dev"

func main() {
	log.Printf("Picup version: %s", version)

	var configPath string
	flag.StringVar(&configPath, "config", "config.ini", "path to config file")
	flag.StringVar(&configPath, "c", "config.ini", "path to config file (shorthand)")
	flag.Parse()

	if flag.NArg() > 0 {
		configPath = flag.Arg(0)
	}

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Error loading config (%s): %v", configPath, err)
	}

	if config.Server.Auth.JWTSecretAutogen {
		log.Printf("WARNING: jwt_secret is not configured and will rotated on app startup")
	}

	// 1. Database Initialization via Factory Pattern
	repo, err := db.NewRepository(config.Server.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize db: %v", err)
	}
	defer repo.Close()

	// init directory
	if err := initializeDirectories(config.GetDirMap()); err != nil {
		log.Fatalf("Error loading config (%s): %v", configPath, err)
	}

	// 2. Prepare Utilities
	exifinsance := exifworker.SetDefault(exifworker.NewExifWorker(config.Server.ExifToolPath))
	defer exifinsance.Close()

	if err := variant.InitProcessors(); err != nil {
		log.Fatalf("Failed to init variant processor: %v", err)
	}

	// 5. Prepare WebServer and handlers
	fileStorage := workflow.NewDefaultStorage(config.Server.TempDir, config.Server.DataDir)
	handler := httpadapter.NewRouter(httpadapter.RouterDependencies{
		Repository:  repo,
		FileStorage: fileStorage,
	})

	// 6. Init Server
	addr := utils.NewBaseURL(config.Server.Host, config.Server.Port, config.Server.DataDir, config.Server.BaseURL)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		log.Printf("🚀 Server starting on http://%s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown listener
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server gracefully...")

}

func initializeDirectories(dirs map[string]string) error {
	for key, dir := range dirs {
		if dir == "" {
			continue
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory %s: %w", key, err)
		}
		log.Printf("\t%16s | %s", key, dir)
	}
	return nil
}
