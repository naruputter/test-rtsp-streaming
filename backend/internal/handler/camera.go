package handler

import (
	"encoding/json"
	"net/http"

	"cctv-backend/internal/config"
	"cctv-backend/internal/stream"
	"github.com/gorilla/mux"
)

type CameraHandler struct {
	manager *stream.Manager
}

func NewCameraHandler(manager *stream.Manager) *CameraHandler {
	return &CameraHandler{manager: manager}
}

func (h *CameraHandler) List(w http.ResponseWriter, r *http.Request) {
	streams := h.manager.List()
	json.NewEncoder(w).Encode(streams)
}

func (h *CameraHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	s, ok := h.manager.Get(id)
	if !ok {
		http.Error(w, "Camera not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(s)
}

func (h *CameraHandler) Start(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	s, ok := h.manager.Get(id)
	if !ok {
		http.Error(w, "Camera not found", http.StatusNotFound)
		return
	}

	if err := h.manager.Start(s.Camera); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *CameraHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.manager.Stop(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *CameraHandler) Add(w http.ResponseWriter, r *http.Request) {
	var cam config.Camera
	if err := json.NewDecoder(r.Body).Decode(&cam); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.manager.AddCamera(cam)
	if cam.Enabled {
		h.manager.Start(cam)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cam)
}

func (h *CameraHandler) Remove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	h.manager.RemoveCamera(id)
	w.WriteHeader(http.StatusNoContent)
}
