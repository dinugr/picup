import { Link, useLocation, useNavigate } from 'react-router-dom';
import { LogOut, Settings, Shield } from 'lucide-react';
import { useAuth } from 'picup/context/AuthContext';

export default function Header() {
  const location = useLocation();
  const navigate = useNavigate();
  const { logoutCurrent } = useAuth();

  const isActive = (path: string) => location.pathname === path || location.pathname.startsWith(path + '/');

  return (
    <header className="sticky top-0 z-50 flex items-center justify-between border-b border-white/10 bg-slate-950/80 px-6 py-4 backdrop-blur-md">
      <div className="flex items-center gap-3.5">
        <Link to="/" className="flex items-center gap-2.5 hover:opacity-80 transition-opacity">
          <h1 className="bg-linear-to-r from-sky-400 via-indigo-400 to-violet-400 bg-clip-text text-xl font-bold tracking-[-0.02em] text-transparent">
            PicUp
          </h1>
        </Link>
        <p className="text-xs text-slate-400">
          Simple Image Service & Storage
        </p>
      </div>

      <nav className="flex items-center gap-6">
        <Link
          to="/"
          className={`text-sm font-medium transition-colors ${
            isActive('/') && !isActive('/image')
              ? 'text-sky-400'
              : 'text-slate-400 hover:text-slate-200'
          }`}
        >
          Gallery
        </Link>
        <Link
          to="/config"
          className={`text-sm font-medium transition-colors ${
            isActive('/config')
              ? 'text-sky-400'
              : 'text-slate-400 hover:text-slate-200'
          }`}
        >
          <Settings className="mr-1 inline" size={14} aria-hidden="true" />
          Config
        </Link>
        <Link
          to="/sessions"
          className={`text-sm font-medium transition-colors ${
            isActive('/sessions')
              ? 'text-sky-400'
              : 'text-slate-400 hover:text-slate-200'
          }`}
        >
          <Shield className="mr-1 inline" size={14} aria-hidden="true" />
          Sessions
        </Link>
        <button
          type="button"
          className="text-sm font-medium text-slate-400 transition-colors hover:text-rose-300"
          onClick={() => {
            void logoutCurrent().finally(() => navigate('/login', { replace: true }));
          }}
        >
          <LogOut className="mr-1 inline" size={14} aria-hidden="true" />
          Log out
        </button>
      </nav>
    </header>
  );
}
