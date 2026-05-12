# CCTV RTSP Streaming Platform

A premium live-streaming dashboard built with **Go** and **Next.js**. It ingests RTSP streams from CCTV cameras using **FFmpeg** and serves them to the browser via **HLS**.

## Features
- **Low-Latency HLS**: Browser-native streaming via `hls.js`.
- **Multi-Camera Dashboard**: Real-time grid view of all cameras.
- **Dynamic Management**: Add/remove cameras through the UI.
- **WebSocket Updates**: Live status changes (Online/Offline/Error).
- **Premium Design**: Dark-mode glassmorphism UI.

## Prerequisites
1. **FFmpeg**: Must be installed and in your system PATH.
   - Windows: `winget install ffmpeg`
   - Mac: `brew install ffmpeg`
   - Linux: `sudo apt install ffmpeg`
2. **Go**: 1.22+
3. **Node.js**: 18+

## How to Run

### 1. Start the Backend
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```
The backend will run on `http://localhost:8080`.

### 2. Start the Frontend
```bash
cd frontend
npm install
npm run dev
```
The frontend will run on `http://localhost:3000`.

## Configuration
Edit `backend/configs/cameras.yaml` to add your real RTSP cameras. A demo Big Buck Bunny RTSP stream is included by default.

```yaml
cameras:
  - id: "front-door"
    name: "Front Door"
    rtsp_url: "rtsp://admin:pass@192.168.1.100:554/stream"
    enabled: true
```

## Troubleshooting
- **No Video?** Ensure FFmpeg is installed. Check the backend terminal for error logs.
- **CORS Issues?** The backend is configured to allow `*`. Ensure you are accessing via `localhost`.
