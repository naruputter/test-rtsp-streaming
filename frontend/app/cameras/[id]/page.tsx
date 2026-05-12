"use client";
import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Stream } from "@/lib/types";
import { api } from "@/lib/api";
import VideoPlayer from "@/components/VideoPlayer";
import StatusBadge from "@/components/StatusBadge";
import { ArrowLeft, RefreshCw, AlertCircle } from "lucide-react";
import Link from "next/link";

export default function CameraDetail() {
  const params = useParams();
  const router = useRouter();
  const id = params.id as string;
  const [stream, setStream] = useState<Stream | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStream = async () => {
      try {
        const data = await api.listCameras();
        const found = data.find(s => s.camera.id === id);
        if (found) setStream(found);
        else router.push("/");
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    fetchStream();
  }, [id, router]);

  if (loading) return <div style={{ padding: '5rem', textAlign: 'center' }}>Loading...</div>;
  if (!stream) return null;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
      <header style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
          <Link href="/" className="btn btn-secondary" style={{ padding: '0.5rem' }}>
            <ArrowLeft size={20} />
          </Link>
          <div>
            <h2 style={{ fontSize: '1.5rem', fontWeight: 700 }}>{stream.camera.name}</h2>
            <p style={{ opacity: 0.5, fontSize: '0.875rem' }}>{id} • {stream.camera.rtsp_url}</p>
          </div>
        </div>
        <StatusBadge status={stream.status} />
      </header>

      <div className="glass" style={{ aspectRatio: '21/9', background: '#000', overflow: 'hidden', minHeight: '500px' }}>
        {stream.status === "running" ? (
          <VideoPlayer url={api.getHLSUrl(id)} />
        ) : (
          <div style={{ width: '100%', height: '100%', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: '1.5rem', color: '#333' }}>
            <AlertCircle size={64} />
            <p style={{ fontSize: '1.125rem' }}>Stream is currently offline</p>
            <button className="btn btn-primary" onClick={() => window.location.reload()}>
              <RefreshCw size={18} />
              Retry Connection
            </button>
          </div>
        )}
      </div>
      
      <section className="glass" style={{ padding: '2rem' }}>
        <h3 style={{ marginBottom: '1rem' }}>Stream Info</h3>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '2rem' }}>
          <div>
            <label style={{ display: 'block', fontSize: '0.75rem', color: '#666', textTransform: 'uppercase', marginBottom: '0.25rem' }}>Protocol</label>
            <span style={{ fontWeight: 600 }}>HLS (HTTP Live Streaming)</span>
          </div>
          <div>
            <label style={{ display: 'block', fontSize: '0.75rem', color: '#666', textTransform: 'uppercase', marginBottom: '0.25rem' }}>Status</label>
            <span style={{ fontWeight: 600, color: stream.status === 'running' ? 'var(--accent)' : 'inherit' }}>{stream.status.toUpperCase()}</span>
          </div>
          <div>
            <label style={{ display: 'block', fontSize: '0.75rem', color: '#666', textTransform: 'uppercase', marginBottom: '0.25rem' }}>Endpoint</label>
            <code style={{ background: 'rgba(255,255,255,0.05)', padding: '0.25rem 0.5rem', borderRadius: '4px', fontSize: '0.875rem' }}>/hls/{id}/index.m3u8</code>
          </div>
        </div>
      </section>
    </div>
  );
}
