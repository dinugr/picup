import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Header from 'picup/components/common/Header';
import ConfirmModal from 'picup/components/common/ConfirmModal';
import SessionList from 'picup/components/auth/SessionList';
import Toast from 'picup/components/common/Toast';
import { fetchSessions, logoutAllSessions, revokeSession } from 'picup/services/api';
import { useAuth } from 'picup/context/AuthContext';
import type { AuthSession } from 'picup/types/api';
import type { ToastState } from 'picup/types/ui';

export default function SessionsPage() {
  const navigate = useNavigate();
  const { logout } = useAuth();
  const [sessions, setSessions] = useState<AuthSession[]>([]);
  const [loading, setLoading] = useState(true);
  const [revokingId, setRevokingId] = useState<string | null>(null);
  const [confirmAll, setConfirmAll] = useState(false);
  const [bulkLoading, setBulkLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [toast, setToast] = useState<ToastState | null>(null);

  const loadSessions = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      setSessions(await fetchSessions());
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : 'Failed to load sessions');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadSessions();
  }, [loadSessions]);

  const handleRevoke = async (session: AuthSession) => {
    try {
      setRevokingId(session.id);
      await revokeSession(session.id);
      if (session.current) {
        logout();
        navigate('/login', { replace: true });
        return;
      }
      setSessions((current) => current.filter((item) => item.id !== session.id));
      setToast({ message: 'Session revoked', type: 'success' });
    } catch (reason) {
      setToast({ message: reason instanceof Error ? reason.message : 'Failed to revoke session', type: 'error' });
    } finally {
      setRevokingId(null);
    }
  };

  const handleLogoutAll = async () => {
    try {
      setBulkLoading(true);
      await logoutAllSessions();
    } catch (reason) {
      setToast({ message: reason instanceof Error ? reason.message : 'Logout all failed', type: 'error' });
    } finally {
      logout();
      setBulkLoading(false);
      setConfirmAll(false);
      navigate('/login', { replace: true });
    }
  };

  return (
    <div className="flex min-h-screen flex-col">
      <Header />
      <main className="mx-auto flex w-full max-w-4xl flex-1 flex-col gap-7 px-6 py-8">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <p className="text-sm font-medium uppercase tracking-[0.18em] text-sky-400">Security</p>
            <h1 className="mt-2 text-3xl font-semibold text-slate-100">Active sessions</h1>
            <p className="mt-2 max-w-xl text-sm text-slate-400">Review signed-in devices and remove access you no longer recognize.</p>
          </div>
          <button type="button" className="btn btn-danger" onClick={() => setConfirmAll(true)} disabled={loading || sessions.length === 0}>
            Log out all devices
          </button>
        </div>

        {error && (
          <div className="glass-panel border-rose-400/30 p-4 text-sm text-rose-200" role="alert">
            <p>{error}</p>
            <button type="button" className="btn btn-secondary mt-3" onClick={() => void loadSessions()}>Try again</button>
          </div>
        )}
        {!error && !loading && sessions.length === 0 && <p className="glass-panel p-6 text-sm text-slate-400">No active sessions found.</p>}
        {!error && (loading || sessions.length > 0) && (
          <SessionList loading={loading} sessions={sessions} revokingId={revokingId} onRevoke={(session) => void handleRevoke(session)} />
        )}
      </main>

      <ConfirmModal
        isOpen={confirmAll}
        title="Log out all devices?"
        message="Every active session, including this device, will be revoked. You will need to sign in again."
        confirmLabel="Log out all devices"
        onConfirm={() => void handleLogoutAll()}
        onCancel={() => setConfirmAll(false)}
        loading={bulkLoading}
      />
      {toast && <Toast message={toast.message} type={toast.type} onClose={() => setToast(null)} />}
    </div>
  );
}
