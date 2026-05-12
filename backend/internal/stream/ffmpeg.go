package stream

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// FFmpegProcess wraps a running FFmpeg subprocess.
type FFmpegProcess struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
	done   chan struct{}
}

// startFFmpeg spawns FFmpeg to convert an RTSP stream to HLS segments
// in outputDir. The returned process owns its own goroutine that waits
// for the subprocess to exit and closes the done channel.
func startFFmpeg(rtspURL, outputDir string) (*FFmpegProcess, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("create hls output dir: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	playlistPath := filepath.Join(outputDir, "index.m3u8")
	segPattern := filepath.Join(outputDir, "seg%05d.ts")

	// FFmpeg flags:
	//  -rtsp_transport tcp  — use TCP for reliability
	//  -c:v copy            — pass through video codec (no re-encode)
	//  -c:a aac             — re-encode audio to AAC for browser compat
	//  -hls_time 2          — 2-second segments
	//  -hls_list_size 10    — keep 10 segments in playlist (rolling)
	//  -hls_flags delete_segments+append_list — delete old .ts files
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-rtsp_transport", "tcp",
		"-i", rtspURL,
		"-c:v", "copy",
		"-c:a", "aac",
		"-ar", "44100",
		"-ac", "2",
		"-f", "hls",
		"-hls_time", "2",
		"-hls_list_size", "10",
		"-hls_flags", "delete_segments+append_list",
		"-hls_segment_filename", segPattern,
		"-y",
		playlistPath,
	)

	// Pipe FFmpeg output to parent logger (prefix per-camera)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("start ffmpeg: %w", err)
	}

	proc := &FFmpegProcess{
		cmd:    cmd,
		cancel: cancel,
		done:   make(chan struct{}),
	}

	// Single goroutine owns Wait() — prevents double-wait races.
	go func() {
		defer close(proc.done)
		_ = cmd.Wait()
	}()

	return proc, nil
}

// Stop cancels the FFmpeg context and blocks until the process exits.
func (p *FFmpegProcess) Stop() {
	p.cancel()
	<-p.done
}

// Done returns a channel that is closed when the process exits.
func (p *FFmpegProcess) Done() <-chan struct{} {
	return p.done
}
