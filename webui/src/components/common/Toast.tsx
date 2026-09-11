import { useEffect, type ReactNode } from 'react';
import { CheckCircle2, AlertCircle, Info, X } from 'lucide-react';
import type { ToastProps, ToastType } from 'picup/types/ui';

export default function Toast({ message, type = 'success', onClose, duration = 4000 }: ToastProps) {
  useEffect(() => {
    if (duration > 0) {
      const timer = setTimeout(() => {
        onClose();
      }, duration);
      return () => clearTimeout(timer);
    }
  }, [duration, onClose]);

  const icons: Record<ToastType, ReactNode> = {
    success: <CheckCircle2 size={18} className="text-emerald-400" />,
    error: <AlertCircle size={18} className="text-rose-400" />,
    info: <Info size={18} className="text-sky-400" />,
  };

  const bgStyles: Record<ToastType, string> = {
    success: 'rgba(52, 211, 153, 0.15)',
    error: 'rgba(251, 113, 133, 0.15)',
    info: 'rgba(56, 189, 248, 0.15)',
  };

  const borderStyles: Record<ToastType, string> = {
    success: 'rgba(52, 211, 153, 0.3)',
    error: 'rgba(251, 113, 133, 0.3)',
    info: 'rgba(56, 189, 248, 0.3)',
  };

  return (
    <div
      className="animate-fade-in"
      style={{
        position: 'fixed',
        bottom: '24px',
        right: '24px',
        zIndex: 1000,
        display: 'flex',
        alignItems: 'center',
        gap: '12px',
        padding: '12px 18px',
        borderRadius: '12px',
        backgroundColor: bgStyles[type],
        border: `1px solid ${borderStyles[type]}`,
        backdropFilter: 'blur(12px)',
        color: '#f8fafc',
        boxShadow: '0 10px 25px rgba(0,0,0,0.5)',
        fontSize: '0.9rem',
        fontWeight: 500,
      }}
    >
      {icons[type]}
      <span>{message}</span>
      <button
        onClick={onClose}
        style={{
          marginLeft: '8px',
          opacity: 0.7,
          display: 'flex',
          alignItems: 'center',
          cursor: 'pointer',
        }}
      >
        <X size={16} />
      </button>
    </div>
  );
}
