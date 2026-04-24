import { BrowserRouter, Routes, Route, Navigate, NavLink } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import LoginPage from './pages/LoginPage';
import EmployeeDashboard from './pages/EmployeeDashboard';
import ManagerDashboard from './pages/ManagerDashboard';

function RequireAuth({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
}

function RequireRole({ role, children }: { role: string; children: React.ReactNode }) {
  const { user } = useAuth();
  if (!user) return <Navigate to="/login" replace />;
  if (user.role !== role) return <Navigate to={user.role === 'manager' ? '/manager' : '/dashboard'} replace />;
  return <>{children}</>;
}

function Header() {
  const { user, logout, isAuthenticated } = useAuth();

  if (!isAuthenticated) return null;

  return (
    <header className="app-header">
      <h1>Leave Management System</h1>
      <nav>
        {user?.role === 'employee' && (
          <NavLink to="/dashboard" className={({ isActive }) => `nav-link${isActive ? ' active' : ''}`}>
            My Dashboard
          </NavLink>
        )}
        {user?.role === 'manager' && (
          <>
            <NavLink to="/manager" className={({ isActive }) => `nav-link${isActive ? ' active' : ''}`}>
              Manager Dashboard
            </NavLink>
            <NavLink to="/dashboard" className={({ isActive }) => `nav-link${isActive ? ' active' : ''}`}>
              My Leaves
            </NavLink>
          </>
        )}
        <button className="btn-logout" onClick={logout}>
          Sign Out
        </button>
      </nav>
    </header>
  );
}

function AppRoutes() {
  const { user } = useAuth();

  return (
    <>
      <Header />
      <Routes>
        <Route path="/login" element={<LoginPage />} />

        <Route
          path="/dashboard"
          element={
            <RequireAuth>
              <EmployeeDashboard />
            </RequireAuth>
          }
        />

        <Route
          path="/manager"
          element={
            <RequireRole role="manager">
              <ManagerDashboard />
            </RequireRole>
          }
        />

        {/* Default redirect */}
        <Route
          path="/"
          element={
            user
              ? <Navigate to={user.role === 'manager' ? '/manager' : '/dashboard'} replace />
              : <Navigate to="/login" replace />
          }
        />

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  );
}
