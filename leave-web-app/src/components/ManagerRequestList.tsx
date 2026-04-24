import React, { useState } from 'react';
import type { LeaveRequest } from '../api/leaveApi';
import { managerApi } from '../api/managerApi';

interface Props {
  requests: LeaveRequest[];
  onUpdate: (req: LeaveRequest) => void;
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

type Action = 'approve' | 'reject';

interface ReviewModal {
  requestId: string;
  action: Action;
  employeeName: string;
}

export default function ManagerRequestList({ requests, onUpdate }: Props) {
  const [modal, setModal] = useState<ReviewModal | null>(null);
  const [comment, setComment] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const openModal = (req: LeaveRequest, action: Action) => {
    setModal({ requestId: req.id, action, employeeName: req.employeeName || req.employeeId });
    setComment('');
    setError('');
  };

  const closeModal = () => {
    setModal(null);
    setComment('');
    setError('');
  };

  const handleReview = async () => {
    if (!modal) return;
    setLoading(true);
    setError('');
    try {
      const fn =
        modal.action === 'approve'
          ? managerApi.approveRequest(modal.requestId, { comment })
          : managerApi.rejectRequest(modal.requestId, { comment });
      const updated = await fn;
      onUpdate(updated);
      closeModal();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Action failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const pending = requests.filter((r) => r.status === 'pending');
  const reviewed = requests.filter((r) => r.status !== 'pending');

  return (
    <>
      {/* Approve/Reject Modal */}
      {modal && (
        <div className="modal-overlay" onClick={closeModal}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h3>
              {modal.action === 'approve' ? 'Approve' : 'Reject'} Request
            </h3>
            <p style={{ color: '#555', marginBottom: '1rem', fontSize: '0.9rem' }}>
              {modal.action === 'approve' ? 'Approving' : 'Rejecting'} leave request for{' '}
              <strong>{modal.employeeName}</strong>.
            </p>
            {error && <div className="alert alert-error">{error}</div>}
            <div className="form-group">
              <label htmlFor="review-comment">Comment (optional)</label>
              <textarea
                id="review-comment"
                value={comment}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setComment(e.target.value)}
                placeholder="Add a comment for the employee…"
                disabled={loading}
              />
            </div>
            <div className="modal-actions">
              <button className="btn btn-secondary" onClick={closeModal} disabled={loading}>
                Cancel
              </button>
              <button
                className={`btn ${modal.action === 'approve' ? 'btn-success' : 'btn-danger'}`}
                onClick={handleReview}
                disabled={loading}
              >
                {loading
                  ? 'Processing…'
                  : modal.action === 'approve'
                  ? 'Approve'
                  : 'Reject'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Pending requests */}
      <div className="card">
        <div className="card-title">Pending Requests ({pending.length})</div>
        {pending.length === 0 ? (
          <div className="empty-state">
            <p>No pending requests to review.</p>
          </div>
        ) : (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Employee</th>
                  <th>Type</th>
                  <th>Start Date</th>
                  <th>End Date</th>
                  <th>Reason</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {pending.map((req) => (
                  <tr key={req.id}>
                    <td>{req.employeeName || req.employeeId}</td>
                    <td>{req.type}</td>
                    <td>{formatDate(req.startDate)}</td>
                    <td>{formatDate(req.endDate)}</td>
                    <td>{req.reason}</td>
                    <td>
                      <div style={{ display: 'flex', gap: '0.5rem' }}>
                        <button
                          className="btn btn-success btn-sm"
                          onClick={() => openModal(req, 'approve')}
                        >
                          Approve
                        </button>
                        <button
                          className="btn btn-danger btn-sm"
                          onClick={() => openModal(req, 'reject')}
                        >
                          Reject
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Reviewed requests */}
      {reviewed.length > 0 && (
        <div className="card">
          <div className="card-title">Reviewed Requests</div>
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Employee</th>
                  <th>Type</th>
                  <th>Start Date</th>
                  <th>End Date</th>
                  <th>Status</th>
                  <th>Comment</th>
                </tr>
              </thead>
              <tbody>
                {reviewed.map((req) => (
                  <tr key={req.id}>
                    <td>{req.employeeName || req.employeeId}</td>
                    <td>{req.type}</td>
                    <td>{formatDate(req.startDate)}</td>
                    <td>{formatDate(req.endDate)}</td>
                    <td>
                      <span className={`badge badge-${req.status}`}>{req.status}</span>
                    </td>
                    <td>{req.comment || '—'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </>
  );
}
