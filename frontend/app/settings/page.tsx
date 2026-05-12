"use client";
import { useState, useEffect } from "react";
import { Camera, Stream } from "@/lib/types";
import { api } from "@/lib/api";
import { Plus, Trash2, ShieldCheck, Server } from "lucide-react";

export default function Settings() {
  const [streams, setStreams] = useState<Stream[]>([]);
  const [newCam, setNewCam] = useState<Partial<Camera>>({ enabled: true });
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    api.listCameras().then(setStreams);
  }, []);

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newCam.id || !newCam.rtsp_url || !newCam.name) return;
    
    setLoading(true);
    try {
      await api.addCamera(newCam as Camera);
      const updated = await api.listCameras();
      setStreams(updated);
      setNewCam({ enabled: true });
    } catch (err) {
      alert("Failed to add camera");
    } finally {
      setLoading(false);
    }
  };

  const handleRemove = async (id: string) => {
    if (!confirm("Are you sure?")) return;
    await api.removeCamera(id);
    setStreams(prev => prev.filter(s => s.Camera.id !== id));
  };

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto' }}>
      <h2 style={{ fontSize: '1.875rem', fontWeight: 800, marginBottom: '2rem' }}>Camera Settings</h2>

      <section className="glass" style={{ padding: '2rem', marginBottom: '3rem' }}>
        <h3 style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '1.5rem' }}>
          <Plus size={20} color="var(--primary)" />
          Add New Camera
        </h3>
        <form onSubmit={handleAdd} style={{ display: 'grid', gap: '1rem' }}>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
            <div className="input-group">
              <label>Unique ID (e.g. cam-01)</label>
              <input 
                type="text" 
                required 
                value={newCam.id || ""} 
                onChange={e => setNewCam({...newCam, id: e.target.value})}
              />
            </div>
            <div className="input-group">
              <label>Friendly Name</label>
              <input 
                type="text" 
                required 
                value={newCam.name || ""} 
                onChange={e => setNewCam({...newCam, name: e.target.value})}
              />
            </div>
          </div>
          <div className="input-group">
            <label>RTSP URL</label>
            <input 
              type="text" 
              required 
              placeholder="rtsp://admin:pass@192.168.1.100:554/stream"
              value={newCam.rtsp_url || ""} 
              onChange={e => setNewCam({...newCam, rtsp_url: e.target.value})}
            />
          </div>
          <button className="btn btn-primary" type="submit" disabled={loading} style={{ marginTop: '1rem' }}>
            {loading ? "Adding..." : "Add Camera"}
          </button>
        </form>
      </section>

      <section>
        <h3 style={{ marginBottom: '1.5rem' }}>Managed Cameras</h3>
        <div style={{ display: 'grid', gap: '1rem' }}>
          {streams.map(s => (
            <div key={s.Camera.id} className="glass" style={{ padding: '1rem 1.5rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                <div style={{ background: 'rgba(255,255,255,0.05)', padding: '0.75rem', borderRadius: '12px' }}>
                  <Server size={20} />
                </div>
                <div>
                  <h4 style={{ fontWeight: 600 }}>{s.Camera.name}</h4>
                  <code style={{ fontSize: '0.75rem', opacity: 0.5 }}>{s.Camera.rtsp_url}</code>
                </div>
              </div>
              <button className="btn btn-danger" onClick={() => handleRemove(s.Camera.id)}>
                <Trash2 size={18} />
              </button>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
