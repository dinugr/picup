import { AlertTriangle, X } from 'lucide-react';
import type { ConfirmModalProps } from 'picup/types/ui';

export default function ConfirmModal({
  isOpen,
  title,
  message,
  onConfirm,
  onCancel,
  loading,
  confirmLabel = 'Confirm',
}: ConfirmModalProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-999 flex items-center justify-center bg-black/75 p-4 backdrop-blur-sm">
      <div className="glass-panel animate-fade-in w-full max-w-110 border border-rose-400/30 bg-slate-900 p-6 shadow-[0_20px_40px_rgba(0,0,0,0.8)]">
        <div className="mb-4 flex items-start justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-rose-500/15 text-rose-400">
              <AlertTriangle size={22} />
            </div>
            <h3 className="text-lg font-semibold">{title || 'Confirm Action'}</h3>
          </div>
          <button onClick={onCancel} className="cursor-pointer text-slate-400 transition hover:text-slate-100">
            <X size={20} />
          </button>
        </div>

        <p className="mb-6 text-sm leading-6 text-slate-400">
          {message}
        </p>

        <div className="flex justify-end gap-3">
          <button className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-white/20 hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-50" onClick={onCancel} disabled={loading}>
            Cancel
          </button>
          <button className="inline-flex cursor-pointer items-center justify-center rounded-lg border border-rose-400/30 bg-rose-500/15 px-4 py-2 text-sm font-medium text-rose-300 transition hover:bg-rose-500/25 disabled:cursor-not-allowed disabled:opacity-50" onClick={onConfirm} disabled={loading}>
            {loading ? 'Working...' : confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
}
