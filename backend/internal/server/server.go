package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"cctv-backend/internal/config"
	"cctv-backend/internal/handler"
	"cctv-backend/internal/stream"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	config  *config.Config
	manager *stream.Manager
	server  *http.Server
}

func NewServer(cfg *config.Config, manager *stream.Manager) *Server {
	r := mux.NewRouter()

	camHandler := handler.NewCameraHandler(manager)
	hlsHandler := handler.NewHLSHandler(cfg.HLSOutputDir)
	wsHandler := handler.NewWSHandler(manager)

	// API Routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/cameras", camHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/cameras", camHandler.Add).Methods(http.MethodPost)
	api.HandleFunc("/cameras/{id}", camHandler.Get).Methods(http.MethodGet)
	api.HandleFunc("/cameras/{id}", camHandler.Remove).Methods(http.MethodDelete)
	api.HandleFunc("/cameras/{id}/start", camHandler.Start).Methods(http.MethodPost)
	api.HandleFunc("/cameras/{id}/stop", camHandler.Stop).Methods(http.MethodPost)

	// HLS Route
	r.PathPrefix("/hls/").Handler(hlsHandler)

	// WebSocket Route
	r.Handle("/ws", wsHandler)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	return &Server{
		config:  cfg,
		manager: manager,
		server: &http.Server{
			Addr:         cfg.ServerAddr,
			Handler:      c.Handler(r),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	log.Printf("[server] listening on %s", s.config.ServerAddr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Printf("[server] shutting down...")
	s.manager.StopAll()
	return s.server.Shutdown(ctx)
}
