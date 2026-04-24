import { useState } from 'react';
import type { LeaveRequest } from '../api/leaveApi';
import { leaveApi } from '../api/leaveApi';

interface Props {
  requests: LeaveRequest[];
  onUpdate: (req: LeaveRequest) => void;
}

function StatusBadge({ status }: { status: LeaveRequest['status'] }) {
  return <span className={`badge badge-${status}`}>{status}</span>;
}

function formatDate(dateStr: string) {
  try {
    return new Date(dateStr).toLocaleDateString(undefined, {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  } catch {
    return dateStr;
  }
}

export default function RequestList({ requests, onUpdate }: Props) {
  const [cancelling, setCancelling] = useState<string | null>(null);
  const [error, setError] = useState('');

  const handleCancel = async (id: string) => {
    if (!window.confirm('Are you sure you want to cancel this request?')) return;
    setCancelling(id);
    setError('');
    try {
      const updated = await leaveApi.cancelRequest(id);
      onUpdate(updated);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to cancel request.');
    } finally {
      setCancelling(null);
    }
  };

  if (requests.length === 0) {
    return (
      <div className="empty-state">
        <p>No leave requests found. Submit a new request above.</p>
      </div>
    );
  }

  return (
    <>
      {error && <div className="alert alert-error">{error}</div>}
      <div className="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>Type</th>
              <th>Start Date</th>
              <th>End Date</th>
              <th>Reason</th>
              <th>Status</th>
              <th>Comment</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {requests.map((req) => (
              <tr key={req.id}>
                <td>{req.type}</td>
                <td>{formatDate(req.startDate)}</td>
                <td>{formatDate(req.endDate)}</td>
                <td>{req.reason}</td>
                <td>
                  <StatusBadge status={req.status} />
                </td>
                <td>{req.comment || '—'}</td>
                <td>
                  {req.status === 'pending' && (
                    <button
                      className="btn btn-danger btn-sm"
                      onClick={() => handleCancel(req.id)}
                      disabled={cancelling === req.id}
                    >
                      {cancelling === req.id ? 'Cancelling…' : 'Cancel'}
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
}
