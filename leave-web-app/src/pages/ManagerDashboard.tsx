import { useEffect, useState } from 'react';
import { managerApi } from '../api/managerApi';
import type { LeaveRequest } from '../api/leaveApi';
import ManagerRequestList from '../components/ManagerRequestList';
import { useAuth } from '../contexts/AuthContext';

export default function ManagerDashboard() {
  const { user } = useAuth();
  const [requests, setRequests] = useState<LeaveRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    managerApi
      .getPendingRequests()
      .then(setRequests)
      .catch((err: Error) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  const handleUpdate = (updated: LeaveRequest) => {
    setRequests((prev) => prev.map((r) => (r.id === updated.id ? updated : r)));
  };

  return (
    <div className="main-content">
      <div className="dashboard-header">
        <h2>Manager Dashboard</h2>
        <span className="user-info">Welcome, {user?.name}</span>
      </div>

      {loading ? (
        <div className="card">
          <div className="loading">
            <div className="spinner" />
            Loading requests…
          </div>
        </div>
      ) : error ? (
        <div className="card">
          <div className="alert alert-error">Failed to load requests: {error}</div>
        </div>
      ) : (
        <ManagerRequestList requests={requests} onUpdate={handleUpdate} />
      )}
    </div>
  );
}
