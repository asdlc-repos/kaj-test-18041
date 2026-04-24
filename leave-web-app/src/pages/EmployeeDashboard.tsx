import { useEffect, useState } from 'react';
import { leaveApi } from '../api/leaveApi';
import type { LeaveRequest } from '../api/leaveApi';
import LeaveBalance from '../components/LeaveBalance';
import LeaveRequestForm from '../components/LeaveRequestForm';
import RequestList from '../components/RequestList';
import { useAuth } from '../contexts/AuthContext';

export default function EmployeeDashboard() {
  const { user } = useAuth();
  const [requests, setRequests] = useState<LeaveRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    leaveApi
      .getMyRequests()
      .then(setRequests)
      .catch((err: Error) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  const handleNewRequest = (req: LeaveRequest) => {
    setRequests((prev) => [req, ...prev]);
  };

  const handleUpdate = (updated: LeaveRequest) => {
    setRequests((prev) => prev.map((r) => (r.id === updated.id ? updated : r)));
  };

  return (
    <div className="main-content">
      <div className="dashboard-header">
        <h2>My Leave Dashboard</h2>
        <span className="user-info">Welcome, {user?.name}</span>
      </div>

      <LeaveBalance />

      <LeaveRequestForm onSuccess={handleNewRequest} />

      <div className="card">
        <div className="card-title">My Requests</div>
        {loading ? (
          <div className="loading">
            <div className="spinner" />
            Loading requests…
          </div>
        ) : error ? (
          <div className="alert alert-error">Failed to load requests: {error}</div>
        ) : (
          <RequestList requests={requests} onUpdate={handleUpdate} />
        )}
      </div>
    </div>
  );
}
