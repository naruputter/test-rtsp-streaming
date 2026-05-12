import { Stream, Camera } from "./types";

const API_BASE = "http://100.105.252.61:8080/api";
const WS_BASE = "ws://100.105.252.61:8080/ws";

export const api = {
  listCameras: async (): Promise<Stream[]> => {
    const res = await fetch(`${API_BASE}/cameras`);
    if (!res.ok) throw new Error("Failed to fetch cameras");
    return res.json();
  },

  startCamera: async (id: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/cameras/${id}/start`, { method: "POST" });
    if (!res.ok) throw new Error("Failed to start camera");
  },

  stopCamera: async (id: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/cameras/${id}/stop`, { method: "POST" });
    if (!res.ok) throw new Error("Failed to stop camera");
  },

  addCamera: async (camera: Camera): Promise<Camera> => {
    const res = await fetch(`${API_BASE}/cameras`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(camera),
    });
    if (!res.ok) throw new Error("Failed to add camera");
    return res.json();
  },

  removeCamera: async (id: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/cameras/${id}`, { method: "DELETE" });
    if (!res.ok) throw new Error("Failed to remove camera");
  },

  getHLSUrl: (id: string): string => {
    return `http://100.105.252.61:8080/hls/${id}/index.m3u8`;
  },

  getWSUrl: (): string => {
    return WS_BASE;
  },
};
