import { Link } from 'react-router-dom';

export default function MovieCard({ movie, delay = 0 }) {
  const genres = {
    'Action': '#E85D24', 'Drama': '#534AB7', 'Sci-Fi': '#185FA5',
    'Comedy': '#0F6E56', 'Horror': '#A32D2D', 'Romance': '#993556',
    'Thriller': '#633806', 'Animation': '#3B6D11',
  };
  const accentColor = genres[movie.genre] || '#534AB7';

  return (
    <Link to={`/movies/${movie.id}`} style={{ textDecoration: 'none' }}>
      <div
        className="fade-up"
        style={{
          animationDelay: `${delay}ms`,
          opacity: 0,
          background: 'var(--ink-2)',
          border: '1px solid var(--ink-4)',
          borderRadius: 'var(--radius-lg)',
          overflow: 'hidden',
          transition: 'transform 0.25s ease, border-color 0.25s ease, box-shadow 0.25s ease',
          cursor: 'pointer',
        }}
        onMouseEnter={e => {
          e.currentTarget.style.transform = 'translateY(-4px)';
          e.currentTarget.style.borderColor = 'var(--smoke)';
          e.currentTarget.style.boxShadow = '0 12px 40px rgba(0,0,0,0.5)';
        }}
        onMouseLeave={e => {
          e.currentTarget.style.transform = 'translateY(0)';
          e.currentTarget.style.borderColor = 'var(--ink-4)';
          e.currentTarget.style.boxShadow = 'none';
        }}
      >
        {/* Poster area */}
        <div style={{
          height: 220,
          background: movie.poster_url
            ? `url(${movie.poster_url}) center/cover no-repeat`
            : `linear-gradient(135deg, ${accentColor}33 0%, var(--ink-3) 100%)`,
          position: 'relative',
          display: 'flex',
          alignItems: 'flex-end',
          padding: '16px',
        }}>
          {!movie.poster_url && (
            <div style={{
              fontFamily: 'var(--font-display)',
              fontSize: 48,
              fontWeight: 300,
              color: `${accentColor}66`,
              position: 'absolute',
              top: '50%',
              left: '50%',
              transform: 'translate(-50%,-50%)',
              letterSpacing: '-0.02em',
            }}>
              {movie.title.charAt(0)}
            </div>
          )}
          <span style={{
            position: 'absolute',
            top: 12,
            right: 12,
            padding: '3px 10px',
            background: 'rgba(0,0,0,0.7)',
            backdropFilter: 'blur(8px)',
            borderRadius: 20,
            fontSize: 11,
            color: 'var(--pearl)',
            letterSpacing: '0.04em',
          }}>
            {movie.duration_min} min
          </span>
        </div>

        {/* Info */}
        <div style={{ padding: '16px 18px 18px' }}>
          <div style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            marginBottom: 6,
          }}>
            <span style={{
              fontSize: 11,
              letterSpacing: '0.08em',
              textTransform: 'uppercase',
              color: accentColor,
              fontWeight: 400,
            }}>
              {movie.genre || 'Film'}
            </span>
            <span style={{ fontSize: 12, color: 'var(--mist)' }}>
              {movie.release_date?.slice(0, 4)}
            </span>
          </div>

          <h3 style={{
            fontFamily: 'var(--font-display)',
            fontSize: 20,
            fontWeight: 400,
            color: 'var(--white)',
            marginBottom: 8,
            lineHeight: 1.3,
          }}>
            {movie.title}
          </h3>

          <p style={{
            fontSize: 13,
            color: 'var(--silver)',
            lineHeight: 1.5,
            display: '-webkit-box',
            WebkitLineClamp: 2,
            WebkitBoxOrient: 'vertical',
            overflow: 'hidden',
          }}>
            {movie.description || 'No description available.'}
          </p>
        </div>
      </div>
    </Link>
  );
}
