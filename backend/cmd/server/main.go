package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cctv-backend/internal/config"
	"cctv-backend/internal/server"
	"cctv-backend/internal/stream"
)

func main() {
	cfg, err := config.Load("configs/cameras.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	manager := stream.NewManager(cfg.HLSOutputDir)

	// Automatically start enabled cameras
	for _, cam := range cfg.Cameras {
		manager.AddCamera(cam)
		if cam.Enabled {
			log.Printf("[main] starting camera %s (%s)", cam.Name, cam.ID)
			if err := manager.Start(cam); err != nil {
				log.Printf("[main] error starting camera %s: %v", cam.ID, err)
			}
		}
	}

	srv := server.NewServer(cfg, manager)

	// Graceful shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("server shutdown failed: %v", err)
		}
	}()

	if err := srv.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
