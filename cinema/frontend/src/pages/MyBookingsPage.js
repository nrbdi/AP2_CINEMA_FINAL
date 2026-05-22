import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../services/api';

function fmt(iso) {
  if (!iso) return '—';
  return new Date(iso).toLocaleString([], {
    month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit',
  });
}

export default function MyBookingsPage() {
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [cancelling, setCancelling] = useState(null);
  const navigate = useNavigate();

  const load = () => {
    api.listBookings()
      .then(data => setBookings(data.bookings || []))
      .catch(() => setBookings([]))
      .finally(() => setLoading(false));
  };

  useEffect(() => { load(); }, []);

  const handleCancel = async (id) => {
    if (!window.confirm('Cancel this booking?')) return;
    setCancelling(id);
    try {
      await api.cancelBooking(id);
      load();
    } catch (e) {
      alert(e.message || 'Could not cancel booking.');
    } finally {
      setCancelling(null);
    }
  };

  return (
    <div className="page">
      <h1 className="page-title fade-up" style={{ opacity: 0 }}>My Bookings</h1>
      <p className="page-sub fade-up fade-up-1">Your cinema reservations and history.</p>

      {loading ? (
        <div className="spinner" />
      ) : bookings.length === 0 ? (
        <div style={{
          textAlign: 'center',
          padding: '80px 0',
          color: 'var(--silver)',
        }}>
          <div style={{
            fontFamily: 'var(--font-display)',
            fontSize: 48,
            fontWeight: 300,
            marginBottom: 12,
            color: 'var(--ink-4)',
          }}>
            No bookings
          </div>
          <p style={{ marginBottom: 24 }}>You haven't made any reservations yet.</p>
          <button className="btn btn-gold" onClick={() => navigate('/')}>
            Browse Films
          </button>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          {bookings.map((b, i) => (
            <div
              key={b.id}
              className="fade-up"
              style={{
                opacity: 0,
                animationDelay: `${i * 60}ms`,
                background: 'var(--ink-2)',
                border: '1px solid var(--ink-4)',
                borderRadius: 'var(--radius-lg)',
                padding: '20px 24px',
                display: 'flex',
                alignItems: 'center',
                gap: 24,
                flexWrap: 'wrap',
              }}
            >
              {/* Status dot */}
              <div style={{
                width: 10, height: 10,
                borderRadius: '50%',
                background: b.status === 'confirmed' ? 'var(--green)' : 'var(--mist)',
                flexShrink: 0,
                boxShadow: b.status === 'confirmed' ? '0 0 8px rgba(75,175,122,0.4)' : 'none',
              }} />

              {/* Main info */}
              <div style={{ flex: 1, minWidth: 200 }}>
                <div style={{
                  fontFamily: 'var(--font-display)',
                  fontSize: 20,
                  fontWeight: 300,
                  marginBottom: 4,
                }}>
                  Booking #{b.id?.slice(0,8).toUpperCase()}
                </div>
                <div style={{ fontSize: 13, color: 'var(--silver)' }}>
                  {b.seat_ids?.length || 0} seat{(b.seat_ids?.length || 0) > 1 ? 's' : ''}
                  <span style={{ color: 'var(--smoke)', margin: '0 8px' }}>·</span>
                  {fmt(b.booked_at)}
                </div>
              </div>

              {/* Price */}
              <div style={{ textAlign: 'right', minWidth: 80 }}>
                <div style={{
                  fontFamily: 'var(--font-display)',
                  fontSize: 24,
                  fontWeight: 300,
                  color: 'var(--gold)',
                }}>
                  ${Number(b.total_price || 0).toFixed(2)}
                </div>
              </div>

              {/* Status badge */}
              <span className={`badge badge-${b.status === 'confirmed' ? 'green' : 'silver'}`}>
                {b.status}
              </span>

              {/* Actions */}
              {b.status === 'confirmed' && (
                <button
                  className="btn btn-danger"
                  style={{ padding: '7px 16px', fontSize: 12 }}
                  onClick={() => handleCancel(b.id)}
                  disabled={cancelling === b.id}
                >
                  {cancelling === b.id ? '…' : 'Cancel'}
                </button>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
