export interface Camera {
  id: string;
  name: string;
  rtsp_url: string;
  enabled: boolean;
}

export type StreamStatus = "running" | "stopped" | "error";

export interface Stream {
  camera: Camera;
  status: StreamStatus;
}

export interface StreamEvent {
  camera_id: string;
  status: StreamStatus;
}
