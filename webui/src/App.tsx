import { Routes, Route } from 'react-router-dom';
import { ConfigProvider } from './context/ConfigContext';
import { AuthProvider } from './context/AuthContext';
import RouteGuard from './components/common/RouteGuard';
import LoginPage from './pages/LoginPage';
import GalleryPage from './pages/GalleryPage';
import ConfigPage from './pages/ConfigPage';
import SessionsPage from './pages/SessionsPage';

/**
 * App Component (Root Router)
 * 
 * Sets up the routing structure and wraps all pages with ConfigProvider.
 * 
 * Routes:
 * - /              : Gallery page (image list)
 * - /image/:id     : Gallery page with detail modal pre-opened
 * - /config        : Configuration display page
 */
export default function App() {
  return (
    <ConfigProvider>
      <AuthProvider>

        <Routes>
          <Route element={<RouteGuard access="guest" />}>
            <Route path="/login" element={<LoginPage />} />
          </Route>

          <Route element={<RouteGuard access="protected" />}>
            <Route path="/" element={<GalleryPage />} />
            <Route path="/image/:id" element={<GalleryPage />} />
            <Route path="/config" element={<ConfigPage />} />
            <Route path="/sessions" element={<SessionsPage />} />
          </Route>
        </Routes>

      </AuthProvider>
    </ConfigProvider>
  );
}
