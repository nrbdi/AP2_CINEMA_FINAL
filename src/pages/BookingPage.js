import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { api } from '../services/api';
import { useAuth } from '../context/AuthContext';
import SeatPicker from '../components/SeatPicker';

const DEMO_SEATS = Array.from({ length: 30 }, (_, i) => ({
  id: `seat-${i+1}`,
  hall_id: 'h1',
  row_number: Math.floor(i / 10) + 1,
  seat_number: (i % 10) + 1,
  seat_type: i < 10 ? 'standard' : i < 20 ? 'premium' : 'vip',
}));

export default function BookingPage() {
  const { showtimeId } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [seats, setSeats] = useState([]);
  const [showtime, setShowtime] = useState(null);
  const [selected, setSelected] = useState([]);
  const [loading, setLoading] = useState(true);
  const [booking, setBooking] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(null);

  useEffect(() => {
    api.getSeats(showtimeId)
      .then(data => setSeats(data.seats || DEMO_SEATS))
      .catch(() => setSeats(DEMO_SEATS))
      .finally(() => setLoading(false));
  }, [showtimeId]);

  const pricePerSeat = showtime?.price || 12.50;
  const total = selected.length * pricePerSeat;

  const handleConfirm = async () => {
    if (selected.length === 0) { setError('Please select at least one seat.'); return; }
    setBooking(true);
    setError('');
    try {
      const data = await api.createBooking({
        showtime_id: showtimeId,
        seat_ids: selected,
      });
      setSuccess(data.booking);
    } catch (e) {
      setError(e.message || 'Booking failed. Please try again.');
    } finally {
      setBooking(false);
    }
  };

  if (success) {
    return (
      <div className="page" style={{ maxWidth: 560, textAlign: 'center', paddingTop: 80 }}>
        <div style={{
          width: 64, height: 64,
          borderRadius: '50%',
          background: 'rgba(75,175,122,0.15)',
          border: '1px solid var(--green)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          margin: '0 auto 24px',
          fontSize: 28,
        }}>
          ✓
        </div>
        <h1 style={{
          fontFamily: 'var(--font-display)',
          fontSize: 40,
          fontWeight: 300,
          marginBottom: 12,
        }}>
          Booking Confirmed
        </h1>
        <p style={{ color: 'var(--silver)', marginBottom: 8 }}>
          Your seats have been reserved successfully.
        </p>
        <p style={{ color: 'var(--silver)', fontSize: 13, marginBottom: 32 }}>
          A confirmation email has been sent to <strong style={{ color: 'var(--white)' }}>{user?.email}</strong>
        </p>
        <div style={{
          background: 'var(--ink-2)',
          border: '1px solid var(--ink-4)',
          borderRadius: 'var(--radius)',
          padding: '20px 24px',
          marginBottom: 32,
          textAlign: 'left',
        }}>
          <div style={{ fontSize: 12, color: 'var(--silver)', marginBottom: 4 }}>Booking ID</div>
          <div style={{ fontFamily: 'var(--font-display)', fontSize: 18, color: 'var(--gold)', marginBottom: 16 }}>
            {success.id?.slice(0,8).toUpperCase()}
          </div>
          <div style={{ fontSize: 12, color: 'var(--silver)', marginBottom: 4 }}>Seats</div>
          <div style={{ color: 'var(--white)', fontSize: 14, marginBottom: 16 }}>
            {selected.length} seat{selected.length > 1 ? 's' : ''}
          </div>
          <div style={{ fontSize: 12, color: 'var(--silver)', marginBottom: 4 }}>Total Paid</div>
          <div style={{ color: 'var(--gold)', fontSize: 22, fontWeight: 400 }}>
            ${Number(success.total_price || total).toFixed(2)}
          </div>
        </div>
        <div style={{ display: 'flex', gap: 12, justifyContent: 'center' }}>
          <button className="btn btn-gold" onClick={() => navigate('/my-bookings')}>
            View My Bookings
          </button>
          <button className="btn btn-ghost" onClick={() => navigate('/')}>
            Browse Films
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="page">
      <button
        onClick={() => navigate(-1)}
        style={{
          background: 'none', border: 'none', color: 'var(--silver)',
          fontSize: 13, cursor: 'pointer', marginBottom: 32,
          display: 'flex', alignItems: 'center', gap: 6,
        }}
      >
        ← Back
      </button>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 320px', gap: 40, alignItems: 'start' }}>
        {/* Seat picker */}
        <div>
          <h1 style={{
            fontFamily: 'var(--font-display)',
            fontSize: 36,
            fontWeight: 300,
            marginBottom: 8,
          }}>
            Select Seats
          </h1>
          <p style={{ color: 'var(--silver)', fontSize: 13, marginBottom: 36 }}>
            Click seats to select or deselect them.
          </p>

          {loading ? <div className="spinner" /> : (
            <SeatPicker seats={seats} onSelectionChange={setSelected} />
          )}
        </div>

        {/* Summary panel */}
        <div style={{ position: 'sticky', top: 88 }}>
          <div style={{
            background: 'var(--ink-2)',
            border: '1px solid var(--ink-4)',
            borderRadius: 'var(--radius-lg)',
            overflow: 'hidden',
          }}>
            <div style={{
              padding: '20px 24px',
              borderBottom: '1px solid var(--ink-4)',
            }}>
              <div style={{ fontSize: 12, color: 'var(--silver)', marginBottom: 4, letterSpacing: '0.06em', textTransform: 'uppercase' }}>
                Order Summary
              </div>
            </div>

            <div style={{ padding: '20px 24px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 12, fontSize: 14 }}>
                <span style={{ color: 'var(--silver)' }}>Seats selected</span>
                <span>{selected.length}</span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 12, fontSize: 14 }}>
                <span style={{ color: 'var(--silver)' }}>Price per seat</span>
                <span>${pricePerSeat.toFixed(2)}</span>
              </div>
              <div style={{ height: 1, background: 'var(--ink-4)', margin: '16px 0' }} />
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 24 }}>
                <span style={{ fontSize: 14 }}>Total</span>
                <span style={{
                  fontFamily: 'var(--font-display)',
                  fontSize: 28,
                  fontWeight: 300,
                  color: 'var(--gold)',
                }}>
                  ${total.toFixed(2)}
                </span>
              </div>

              {error && (
                <div style={{
                  padding: '10px 14px',
                  background: 'rgba(200,75,75,0.1)',
                  border: '1px solid rgba(200,75,75,0.3)',
                  borderRadius: 'var(--radius-sm)',
                  fontSize: 13,
                  color: 'var(--red)',
                  marginBottom: 16,
                }}>
                  {error}
                </div>
              )}

              <button
                className="btn btn-gold"
                style={{ width: '100%', justifyContent: 'center', padding: '13px' }}
                onClick={handleConfirm}
                disabled={booking || selected.length === 0}
              >
                {booking ? 'Processing…' : 'Confirm Booking'}
              </button>

              <p style={{
                fontSize: 11,
                color: 'var(--mist)',
                textAlign: 'center',
                marginTop: 12,
                lineHeight: 1.5,
              }}>
                By confirming you agree to our terms. A receipt will be emailed to you.
              </p>
            </div>
          </div>
        </div>
      </div>

      <style>{`@media(max-width:700px){ .book-grid { grid-template-columns: 1fr !important; } }`}</style>
    </div>
  );
}
