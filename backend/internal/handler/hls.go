package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type HLSHandler struct {
	outputDir string
}

func NewHLSHandler(outputDir string) *HLSHandler {
	return &HLSHandler{outputDir: outputDir}
}

func (h *HLSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// The URL expected is /hls/{cameraID}/{filename}
	path := strings.TrimPrefix(r.URL.Path, "/hls/")
	fullPath := filepath.Join(h.outputDir, path)

	// Basic security: prevent directory traversal
	if !strings.HasPrefix(filepath.Clean(fullPath), filepath.Clean(h.outputDir)) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	// Set appropriate headers for HLS
	if strings.HasSuffix(fullPath, ".m3u8") {
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	} else if strings.HasSuffix(fullPath, ".ts") {
		w.Header().Set("Content-Type", "video/MP2T")
	}

	// Allow CORS for frontend
	w.Header().Set("Access-Control-Allow-Origin", "*")

	http.ServeFile(w, r, fullPath)
}
