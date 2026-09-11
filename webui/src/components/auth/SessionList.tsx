import LoadingState from 'picup/components/common/LoadingState';
import { Monitor, Smartphone, Tablet, Trash2 } from 'lucide-react';
import type { AuthSession } from 'picup/types/api';

type SessionListProps = {
  sessions: AuthSession[];
  loading: boolean;
  revokingId: string | null;
  onRevoke: (session: AuthSession) => void;
};

function deviceIcon(userAgent: string) {
  const value = userAgent.toLowerCase();
  if (value.includes('mobile') || value.includes('android') || value.includes('iphone')) return Smartphone;
  if (value.includes('tablet') || value.includes('ipad')) return Tablet;
  return Monitor;
}

function formatDate(value?: string): string {
  if (!value) return 'No activity recorded';
  return new Intl.DateTimeFormat(undefined, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value));
}

export default function SessionList({ sessions, loading, revokingId, onRevoke }: SessionListProps) {
  if (loading) {
    return <LoadingState message="Loading active sessions..." />;
  }

  return (
    <div className="space-y-3" aria-label="Active sessions">
      {sessions.map((session) => {
        const DeviceIcon = deviceIcon(session.userAgent);
        return (
          <article key={session.id} className="glass-panel flex flex-col gap-4 p-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex min-w-0 items-start gap-3">
              <div className="rounded-lg bg-sky-400/10 p-2 text-sky-300" aria-hidden="true">
                <DeviceIcon size={20} />
              </div>
              <div className="min-w-0">
                <div className="flex flex-wrap items-center gap-2">
                  <h2 className="font-medium text-slate-100">{session.userAgent || 'Unknown device'}</h2>
                  {session.current && <span className="badge badge-blue">Current device</span>}
                </div>
                <p className="mt-1 text-sm text-slate-400">Signed in {formatDate(session.createdAt)}</p>
                <p className="text-sm text-slate-500">Last active {formatDate(session.lastSeenAt)}</p>
              </div>
            </div>
            <button
              type="button"
              className="btn btn-danger shrink-0 self-end sm:self-auto"
              onClick={() => onRevoke(session)}
              disabled={revokingId === session.id}
              aria-label={`Revoke ${session.current ? 'current' : ''} session`}
            >
              <Trash2 size={16} aria-hidden="true" />
              {revokingId === session.id ? 'Revoking...' : 'Revoke'}
            </button>
          </article>
        );
      })}
    </div>
  );
}
