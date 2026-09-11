import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useAuth } from 'picup/context/AuthContext';

type RouteGuardProps = {
  access: 'protected' | 'guest';
};

export default function RouteGuard({ access }: RouteGuardProps) {
  const location = useLocation();
  const { isAuthenticated, isAuthLoading } = useAuth();

  if (isAuthLoading) {
    return null;
  }

  if (access === 'protected' && !isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  if (access === 'guest' && isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return <Outlet />;
}