import { useEffect, useState } from 'react';
import { leaveApi } from '../api/leaveApi';
import type { LeaveBalance as LeaveBalanceType } from '../api/leaveApi';

export default function LeaveBalance() {
  const [balance, setBalance] = useState<LeaveBalanceType | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    leaveApi
      .getBalance()
      .then(setBalance)
      .catch((err: Error) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading)
    return (
      <div className="loading">
        <div className="spinner" />
        Loading balance…
      </div>
    );

  if (error)
    return <div className="alert alert-error">Failed to load balance: {error}</div>;

  if (!balance) return null;

  const entries = Object.entries(balance).filter(([, v]) => typeof v === 'number');

  return (
    <div className="card">
      <div className="card-title">Leave Balance</div>
      <div className="balance-grid">
        {entries.map(([type, days]) => (
          <div key={type} className="balance-item">
            <div className="balance-value">{days}</div>
            <div className="balance-label">{type.charAt(0).toUpperCase() + type.slice(1)}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
