import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { api } from '../services/api';
import { useAuth } from '../context/AuthContext';

const DEMO_SHOWTIMES = [
  { id: 'st1', movie_id: '1', hall_id: 'h1', start_time: new Date(Date.now() + 3600000*2).toISOString(), end_time: new Date(Date.now() + 3600000*5).toISOString(), price: 12.50 },
  { id: 'st2', movie_id: '1', hall_id: 'h2', start_time: new Date(Date.now() + 3600000*6).toISOString(), end_time: new Date(Date.now() + 3600000*9).toISOString(), price: 15.00 },
  { id: 'st3', movie_id: '1', hall_id: 'h1', start_time: new Date(Date.now() + 3600000*24).toISOString(), end_time: new Date(Date.now() + 3600000*27).toISOString(), price: 12.50 },
];

const DEMO_MOVIE = { id: '1', title: 'Demo Film', genre: 'Sci-Fi', duration_min: 166, release_date: '2024-03-01', description: 'This is a demo movie. Connect your API to see real content.' };

function fmt(iso) {
  if (!iso) return '';
  const d = new Date(iso);
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}
function fmtDate(iso) {
  if (!iso) return '';
  return new Date(iso).toLocaleDateString([], { weekday: 'short', month: 'short', day: 'numeric' });
}

export default function MovieDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [movie, setMovie] = useState(null);
  const [showtimes, setShowtimes] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      api.getMovie(id).catch(() => DEMO_MOVIE),
      api.listShowtimes(id).catch(() => ({ showtimes: DEMO_SHOWTIMES })),
    ]).then(([movieData, stData]) => {
      setMovie(movieData.movie || DEMO_MOVIE);
      setShowtimes(stData.showtimes || DEMO_SHOWTIMES);
    }).finally(() => setLoading(false));
  }, [id]);

  const handleBook = (showtimeId) => {
    if (!user) { navigate('/login'); return; }
    navigate(`/book/${showtimeId}`);
  };

  // Group showtimes by date
  const byDate = {};
  showtimes.forEach(st => {
    const day = new Date(st.start_time).toDateString();
    if (!byDate[day]) byDate[day] = [];
    byDate[day].push(st);
  });

  if (loading) return <div className="spinner" />;
  if (!movie) return null;

  return (
    <div>
      {/* Hero banner */}
      <div style={{
        background: movie.poster_url
          ? `linear-gradient(to bottom, rgba(10,10,11,0) 0%, rgba(10,10,11,1) 100%), url(${movie.poster_url}) center/cover no-repeat`
          : 'linear-gradient(135deg, var(--ink-3) 0%, var(--ink-2) 100%)',
        minHeight: 320,
        display: 'flex',
        alignItems: 'flex-end',
        padding: '0 0 48px',
      }}>
        <div style={{ maxWidth: 1200, margin: '0 auto', padding: '0 24px', width: '100%' }}>
          <button
            onClick={() => navigate(-1)}
            style={{
              background: 'none', border: 'none', color: 'var(--silver)',
              fontSize: 13, cursor: 'pointer', marginBottom: 24,
              display: 'flex', alignItems: 'center', gap: 6,
              letterSpacing: '0.04em',
            }}
          >
            ← Back
          </button>
          <div style={{ display: 'flex', alignItems: 'flex-end', gap: 24, flexWrap: 'wrap' }}>
            <div style={{ flex: 1, minWidth: 240 }}>
              <div style={{ marginBottom: 10 }}>
                <span style={{
                  fontSize: 11, letterSpacing: '0.1em', textTransform: 'uppercase',
                  color: 'var(--gold)',
                }}>
                  {movie.genre}
                </span>
                <span style={{ color: 'var(--smoke)', margin: '0 10px' }}>·</span>
                <span style={{ fontSize: 12, color: 'var(--silver)' }}>{movie.duration_min} min</span>
                <span style={{ color: 'var(--smoke)', margin: '0 10px' }}>·</span>
                <span style={{ fontSize: 12, color: 'var(--silver)' }}>{movie.release_date?.slice(0,4)}</span>
              </div>
              <h1 style={{
                fontFamily: 'var(--font-display)',
                fontSize: 'clamp(28px, 4vw, 48px)',
                fontWeight: 300,
                marginBottom: 12,
                lineHeight: 1.1,
              }}>
                {movie.title}
              </h1>
              <p style={{ color: 'var(--silver)', maxWidth: 560, lineHeight: 1.7, fontSize: 14 }}>
                {movie.description}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Showtimes */}
      <div className="page" style={{ paddingTop: 40 }}>
        <h2 style={{
          fontFamily: 'var(--font-display)',
          fontSize: 28,
          fontWeight: 300,
          marginBottom: 32,
        }}>
          Showtimes
        </h2>

        {Object.keys(byDate).length === 0 ? (
          <div style={{ color: 'var(--silver)', fontSize: 14 }}>No showtimes available.</div>
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', gap: 32 }}>
            {Object.entries(byDate).map(([day, sts]) => (
              <div key={day}>
                <div style={{
                  fontSize: 12,
                  letterSpacing: '0.1em',
                  textTransform: 'uppercase',
                  color: 'var(--silver)',
                  marginBottom: 14,
                  display: 'flex',
                  alignItems: 'center',
                  gap: 12,
                }}>
                  {fmtDate(sts[0].start_time)}
                  <div style={{ flex: 1, height: 1, background: 'var(--ink-4)' }} />
                </div>

                <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
                  {sts.map(st => (
                    <div
                      key={st.id}
                      style={{
                        background: 'var(--ink-2)',
                        border: '1px solid var(--ink-4)',
                        borderRadius: 'var(--radius)',
                        padding: '16px 20px',
                        minWidth: 160,
                        transition: 'border-color 0.2s, box-shadow 0.2s',
                      }}
                      onMouseEnter={e => {
                        e.currentTarget.style.borderColor = 'var(--smoke)';
                        e.currentTarget.style.boxShadow = '0 4px 20px rgba(0,0,0,0.3)';
                      }}
                      onMouseLeave={e => {
                        e.currentTarget.style.borderColor = 'var(--ink-4)';
                        e.currentTarget.style.boxShadow = 'none';
                      }}
                    >
                      <div style={{
                        fontFamily: 'var(--font-display)',
                        fontSize: 28,
                        fontWeight: 300,
                        color: 'var(--white)',
                        marginBottom: 2,
                      }}>
                        {fmt(st.start_time)}
                      </div>
                      <div style={{ fontSize: 12, color: 'var(--silver)', marginBottom: 16 }}>
                        ends {fmt(st.end_time)}
                      </div>
                      <div style={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        gap: 12,
                      }}>
                        <span style={{ color: 'var(--gold)', fontSize: 16, fontWeight: 400 }}>
                          ${Number(st.price).toFixed(2)}
                        </span>
                        <button
                          className="btn btn-gold"
                          style={{ padding: '7px 16px', fontSize: 12 }}
                          onClick={() => handleBook(st.id)}
                        >
                          Book
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
