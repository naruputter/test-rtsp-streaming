import { StreamStatus } from "@/lib/types";

export default function StatusBadge({ status }: { status: StreamStatus }) {
  const styles: Record<StreamStatus, { bg: string; color: string; label: string }> = {
    running: { bg: "rgba(16, 185, 129, 0.1)", color: "#10b981", label: "Live" },
    stopped: { bg: "rgba(255, 255, 255, 0.05)", color: "#9ca3af", label: "Offline" },
    error: { bg: "rgba(239, 68, 68, 0.1)", color: "#ef4444", label: "Error" },
  };

  const { bg, color, label } = styles[status];

  return (
    <div style={{ 
      display: 'inline-flex', 
      alignItems: 'center', 
      gap: '0.5rem', 
      padding: '0.25rem 0.75rem', 
      borderRadius: '9999px',
      fontSize: '0.75rem',
      fontWeight: 600,
      textTransform: 'uppercase',
      background: bg,
      color: color,
      border: `1px solid ${color}33`
    }}>
      {status === 'running' && (
        <span style={{ 
          width: '8px', 
          height: '8px', 
          borderRadius: '50%', 
          background: color,
          boxShadow: `0 0 8px ${color}`,
          animation: 'pulse 2s infinite'
        }} />
      )}
      {label}
    </div>
  );
}
