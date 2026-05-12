"use client";
import { Stream } from "@/lib/types";
import { api } from "@/lib/api";
import StatusBadge from "./StatusBadge";
import VideoPlayer from "./VideoPlayer";
import { Play, Square, AlertCircle, Maximize2 } from "lucide-react";
import { useState } from "react";
import Link from "next/link";

export default function CameraCard({ stream: initialStream }: { stream: Stream }) {
  const [stream, setStream] = useState(initialStream);
  const [isLoading, setIsLoading] = useState(false);

  const toggleStream = async () => {
    setIsLoading(true);
    try {
      if (stream.status === "running") {
        await api.stopCamera(stream.camera.id);
        setStream({ ...stream, status: "stopped" });
      } else {
        await api.startCamera(stream.camera.id);
        setStream({ ...stream, status: "running" });
      }
    } catch (err) {
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="glass" style={{ display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
      <div style={{ padding: '1rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <h3 style={{ fontWeight: 600 }}>{stream.camera.name}</h3>
          <p style={{ fontSize: '0.75rem', color: '#666' }}>{stream.camera.id}</p>
        </div>
        <StatusBadge status={stream.status} />
      </div>

      <div style={{ aspectRatio: '16/9', background: '#000', position: 'relative' }}>
        {stream.status === "running" ? (
          <VideoPlayer url={api.getHLSUrl(stream.camera.id)} />
        ) : (
          <div style={{ width: '100%', height: '100%', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: '1rem', color: '#444' }}>
            {stream.status === "error" ? <AlertCircle size={48} color="var(--error)" /> : <Square size={48} />}
            <span style={{ fontSize: '0.875rem' }}>
              {stream.status === "error" ? "Connection Failed" : "Stream Offline"}
            </span>
          </div>
        )}
      </div>

      <div style={{ padding: '1rem', display: 'flex', gap: '0.75rem' }}>
        <button 
          className={`btn ${stream.status === "running" ? 'btn-danger' : 'btn-primary'}`}
          onClick={toggleStream}
          disabled={isLoading}
          style={{ flex: 1, justifyContent: 'center' }}
        >
          {stream.status === "running" ? <Square size={18} /> : <Play size={18} />}
          {stream.status === "running" ? "Stop Stream" : "Start Stream"}
        </button>
        
        <Link href={`/cameras/${stream.camera.id}`} className="btn btn-secondary" title="View Fullscreen">
          <Maximize2 size={18} />
        </Link>
      </div>
    </div>
  );
}
