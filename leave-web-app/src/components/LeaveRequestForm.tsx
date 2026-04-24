import React, { useState } from 'react';
import { leaveApi } from '../api/leaveApi';
import type { LeaveRequest } from '../api/leaveApi';

interface Props {
  onSuccess: (req: LeaveRequest) => void;
}

const LEAVE_TYPES = ['Annual', 'Sick', 'Casual', 'Maternity', 'Paternity', 'Unpaid'];

export default function LeaveRequestForm({ onSuccess }: Props) {
  const [type, setType] = useState(LEAVE_TYPES[0]);
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [reason, setReason] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const today = new Date().toISOString().split('T')[0];

  const validate = (): string | null => {
    if (!startDate) return 'Start date is required.';
    if (!endDate) return 'End date is required.';
    if (endDate < startDate) return 'End date must be on or after the start date.';
    if (!reason.trim()) return 'Reason is required.';
    return null;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const validationError = validate();
    if (validationError) {
      setError(validationError);
      return;
    }

    setLoading(true);
    setError('');
    setSuccess('');

    try {
      const req = await leaveApi.createRequest({ type, startDate, endDate, reason: reason.trim() });
      setSuccess('Leave request submitted successfully!');
      setStartDate('');
      setEndDate('');
      setReason('');
      setType(LEAVE_TYPES[0]);
      onSuccess(req);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to submit request.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="card">
      <div className="card-title">New Leave Request</div>

      {error && <div className="alert alert-error">{error}</div>}
      {success && <div className="alert alert-success">{success}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="leaveType">Leave Type</label>
          <select id="leaveType" value={type} onChange={(e) => setType(e.target.value)} disabled={loading}>
            {LEAVE_TYPES.map((t) => (
              <option key={t} value={t}>
                {t}
              </option>
            ))}
          </select>
        </div>

        <div className="form-row">
          <div className="form-group">
            <label htmlFor="startDate">Start Date</label>
            <input
              id="startDate"
              type="date"
              value={startDate}
              min={today}
              onChange={(e) => setStartDate(e.target.value)}
              disabled={loading}
            />
          </div>
          <div className="form-group">
            <label htmlFor="endDate">End Date</label>
            <input
              id="endDate"
              type="date"
              value={endDate}
              min={startDate || today}
              onChange={(e) => setEndDate(e.target.value)}
              disabled={loading}
            />
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="reason">Reason</label>
          <textarea
            id="reason"
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            placeholder="Briefly describe the reason for your leave…"
            disabled={loading}
          />
        </div>

        <button type="submit" className="btn btn-primary" disabled={loading}>
          {loading ? 'Submitting…' : 'Submit Request'}
        </button>
      </form>
    </div>
  );
}
