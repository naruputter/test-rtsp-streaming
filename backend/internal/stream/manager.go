package stream

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"cctv-backend/internal/config"
)

// Status represents the lifecycle state of a stream.
type Status string

const (
	StatusStopped Status = "stopped"
	StatusRunning Status = "running"
	StatusError   Status = "error"
)

// Stream holds the state of one camera's stream.
type Stream struct {
	Camera  config.Camera `json:"camera"`
	Status  Status        `json:"status"`
	process *FFmpegProcess
}

// Event is emitted on the Manager's event channel whenever a stream
// changes state.
type Event struct {
	CameraID string `json:"camera_id"`
	Status   Status `json:"status"`
}

// Manager tracks all cameras and their running FFmpeg processes.
type Manager struct {
	mu           sync.RWMutex
	streams      map[string]*Stream // keyed by camera ID
	hlsOutputDir string
	events       chan Event
}

// NewManager creates a Manager that writes HLS output under hlsOutputDir.
func NewManager(hlsOutputDir string) *Manager {
	return &Manager{
		streams:      make(map[string]*Stream),
		hlsOutputDir: hlsOutputDir,
		events:       make(chan Event, 256),
	}
}

// Events returns a read-only channel of stream state changes.
func (m *Manager) Events() <-chan Event {
	return m.events
}

// AddCamera registers a camera without starting its stream.
// Idempotent — calling twice with the same ID is a no-op.
func (m *Manager) AddCamera(cam config.Camera) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.streams[cam.ID]; !exists {
		m.streams[cam.ID] = &Stream{Camera: cam, Status: StatusStopped}
	}
}

// RemoveCamera stops and removes a camera from the manager.
func (m *Manager) RemoveCamera(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.streams[id]
	if !ok {
		return
	}
	if s.process != nil {
		s.process.Stop()
	}
	delete(m.streams, id)
}

// Start launches an FFmpeg process for the given camera.
// The camera must have been registered via AddCamera first.
func (m *Manager) Start(cam config.Camera) error {
	// Ensure camera is registered.
	m.AddCamera(cam)

	m.mu.Lock()
	defer m.mu.Unlock()

	s := m.streams[cam.ID]
	if s.Status == StatusRunning {
		return fmt.Errorf("stream %s is already running", cam.ID)
	}

	outputDir := filepath.Join(m.hlsOutputDir, cam.ID)
	proc, err := startFFmpeg(cam.RTSPURL, outputDir)
	if err != nil {
		s.Status = StatusError
		m.emit(cam.ID, StatusError)
		return err
	}

	s.process = proc
	s.Status = StatusRunning
	m.emit(cam.ID, StatusRunning)

	// Watch for unexpected process exit.
	go func() {
		<-proc.Done()

		m.mu.Lock()
		if cur, ok := m.streams[cam.ID]; ok && cur.process == proc {
			cur.Status = StatusStopped
			cur.process = nil
		}
		m.mu.Unlock()

		log.Printf("[stream] camera %s exited", cam.ID)
		m.emit(cam.ID, StatusStopped)
	}()

	return nil
}

// Stop terminates the running stream for the given camera ID.
func (m *Manager) Stop(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.streams[id]
	if !ok {
		return fmt.Errorf("camera %s not found", id)
	}
	if s.process == nil {
		return nil // already stopped
	}

	proc := s.process
	s.process = nil
	s.Status = StatusStopped

	// Release lock before blocking on Stop().
	m.mu.Unlock()
	proc.Stop()
	m.mu.Lock()

	m.emit(id, StatusStopped)
	return nil
}

// StopAll terminates every running stream. Called on graceful shutdown.
func (m *Manager) StopAll() {
	m.mu.Lock()
	procs := make([]*FFmpegProcess, 0)
	for _, s := range m.streams {
		if s.process != nil {
			procs = append(procs, s.process)
			s.process = nil
			s.Status = StatusStopped
		}
	}
	m.mu.Unlock()

	for _, p := range procs {
		p.Stop()
	}
}

// List returns a snapshot of all streams.
func (m *Manager) List() []Stream {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]Stream, 0, len(m.streams))
	for _, s := range m.streams {
		out = append(out, *s)
	}
	return out
}

// Get returns the stream state for a single camera.
func (m *Manager) Get(id string) (Stream, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.streams[id]
	if !ok {
		return Stream{}, false
	}
	return *s, true
}

// UpdateCamera replaces the camera metadata (name, URL) for an existing entry.
func (m *Manager) UpdateCamera(cam config.Camera) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.streams[cam.ID]
	if !ok {
		return fmt.Errorf("camera %s not found", cam.ID)
	}
	s.Camera = cam
	return nil
}

func (m *Manager) emit(cameraID string, status Status) {
	select {
	case m.events <- Event{CameraID: cameraID, Status: status}:
	default: // drop if buffer full
	}
}
