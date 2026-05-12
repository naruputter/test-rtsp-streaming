"use client";
import { useEffect, useState } from "react";
import { Stream, StreamEvent } from "@/lib/types";
import { api } from "@/lib/api";
import CameraCard from "@/components/CameraCard";
import { RefreshCw, LayoutGrid } from "lucide-react";

export default function Dashboard() {
  const [streams, setStreams] = useState<Stream[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStreams = async () => {
    try {
      setLoading(true);
      const data = await api.listCameras();
      setStreams(data);
      setError(null);
    } catch (err) {
      setError("Failed to connect to backend server");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStreams();

    // WebSocket for live updates
    const ws = new WebSocket(api.getWSUrl());
    ws.onmessage = (event) => {
      const update: StreamEvent = JSON.parse(event.data);
      setStreams((prev) =>
        prev.map((s) => s.camera.id === update.camera_id ? { ...s, status: update.status } : s)
      );
    };

    return () => ws.close();
  }, []);

  return (
    <div>
      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
        <div>
          <h2 style={{ fontSize: '1.875rem', fontWeight: 800 }}>Live Dashboard</h2>
          <p style={{ color: '#666' }}>Monitoring {streams.length} cameras</p>
        </div>
        <button
          className="btn btn-secondary"
          onClick={fetchStreams}
          disabled={loading}
        >
          <RefreshCw size={18} className={loading ? 'spin' : ''} />
          Refresh
        </button>
      </header>

      {error && (
        <div className="glass" style={{ padding: '2rem', textAlign: 'center', border: '1px solid var(--error)', color: 'var(--error)', marginBottom: '2rem' }}>
          {error}
        </div>
      )}

      {loading && streams.length === 0 ? (
        <div style={{ display: 'flex', justifyContent: 'center', padding: '5rem' }}>
          <div className="spin" style={{ width: '40px', height: '40px', border: '4px solid var(--primary)', borderTopColor: 'transparent', borderRadius: '50%' }} />
        </div>
      ) : (
        <div className="grid">
          {streams.map((s) => (
            <CameraCard key={s.camera.id} stream={s} />
          ))}
          {streams.length === 0 && !loading && (
            <div className="glass" style={{ gridColumn: '1 / -1', padding: '4rem', textAlign: 'center', color: '#666' }}>
              <LayoutGrid size={48} style={{ margin: '0 auto 1rem', opacity: 0.2 }} />
              <p>No cameras found. Go to Settings to add one.</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
