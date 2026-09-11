import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from 'react';
import { logoutAllSessions, logoutCurrentSession } from 'picup/services/api';

export type AuthState = {
  accessToken: string | null;
};

export type AuthContextValue = {
  auth: AuthState;
  isAuthenticated: boolean;
  isAuthLoading: boolean;
  login: (accessToken: string) => void;
  logout: () => void;
  logoutCurrent: () => Promise<void>;
  logoutAll: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

// With HttpOnly cookie auth, frontend cannot read the JWT value.
// We keep only an in-memory flag to drive route guards.

export function AuthProvider({ children }: { children: ReactNode }) {
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [isAuthLoading, setIsAuthLoading] = useState<boolean>(true);

  useEffect(() => {
    const handleUnauthorized = () => setAccessToken(null);
    window.addEventListener('picup:unauthorized', handleUnauthorized);

    return () => window.removeEventListener('picup:unauthorized', handleUnauthorized);
  }, []);

  useEffect(() => {
    // With HttpOnly cookie auth, frontend cannot read the JWT value.
    // Instead, we validate the cookie via a lightweight endpoint.
    const rehydrate = async () => {
      try {
        const apiBase = import.meta.env.VITE_API_BASE_URL || '';
        const res = await fetch(`${apiBase}/api/auth/session`, {
          method: 'GET',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
        });

        if (res.ok) {
          setAccessToken('cookie');
        }
      } catch {
        // Route guards handle unavailable authentication services.
      } finally {
        setIsAuthLoading(false);
      }
    };

    void rehydrate();
  }, []);

  const login = (token: string) => {
    setAccessToken(token);
  };

  const logout = () => {
    setAccessToken(null);
  };

  const logoutCurrent = async () => {
    try {
      await logoutCurrentSession();
    } finally {
      logout();
    }
  };

  const logoutAll = async () => {
    try {
      await logoutAllSessions();
    } finally {
      logout();
    }
  };

  const value = useMemo<AuthContextValue>(
    () => ({
      auth: { accessToken },
      isAuthenticated: !!accessToken,
      isAuthLoading,
      login,
      logout,
      logoutCurrent,
      logoutAll,
    }),
    [accessToken, isAuthLoading],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return ctx;
}
