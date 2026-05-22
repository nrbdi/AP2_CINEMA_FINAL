import { useState, useEffect } from 'react';
import { api } from '../services/api';
import MovieCard from '../components/MovieCard';

const GENRES = ['All', 'Action', 'Drama', 'Sci-Fi', 'Comedy', 'Horror', 'Romance', 'Thriller', 'Animation'];

// Fallback demo movies when API isn't running
const DEMO_MOVIES = [
  { id: '1', title: 'Dune: Part Two', genre: 'Sci-Fi', duration_min: 166, release_date: '2024-03-01', description: 'Paul Atreides unites with Chani and the Fremen while seeking revenge against the conspirators who destroyed his family.' },
  { id: '2', title: 'Oppenheimer', genre: 'Drama', duration_min: 180, release_date: '2023-07-21', description: 'The story of American scientist J. Robert Oppenheimer and his role in the development of the atomic bomb.' },
  { id: '3', title: 'The Dark Knight', genre: 'Action', duration_min: 152, release_date: '2008-07-18', description: 'When the menace known as the Joker wreaks havoc and chaos on Gotham, Batman must accept one of the greatest tests.' },
  { id: '4', title: 'Interstellar', genre: 'Sci-Fi', duration_min: 169, release_date: '2014-11-07', description: 'A team of explorers travel through a wormhole in space in an attempt to ensure humanity\'s survival.' },
  { id: '5', title: 'Parasite', genre: 'Thriller', duration_min: 132, release_date: '2019-10-11', description: 'All four members of a Ki-taek\'s family are unemployed, living in a semi-basement and looking for opportunities.' },
  { id: '6', title: 'Past Lives', genre: 'Romance', duration_min: 106, release_date: '2023-06-02', description: 'Two childhood sweethearts are separated when one emigrates from South Korea. Decades later, they reunite in New York.' },
  { id: '7', title: 'Poor Things', genre: 'Comedy', duration_min: 141, release_date: '2023-12-08', description: 'The incredible tale of Bella Baxter, a young woman brought back to life by an eccentric scientist.' },
  { id: '8', title: 'Alien: Romulus', genre: 'Horror', duration_min: 119, release_date: '2024-08-16', description: 'A group of young space colonizers come face to face with the most terrifying life form in the universe.' },
];

export default function MoviesPage() {
  const [movies, setMovies] = useState([]);
  const [genre, setGenre] = useState('All');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    setLoading(true);
    api.listMovies(genre === 'All' ? '' : genre)
      .then(data => setMovies(data.movies || []))
      .catch(() => {
        setMovies(genre === 'All' ? DEMO_MOVIES : DEMO_MOVIES.filter(m => m.genre === genre));
        setError('Showing demo content — API not connected');
      })
      .finally(() => setLoading(false));
  }, [genre]);

  const filtered = movies;

  return (
    <div className="page">
      {/* Hero */}
      <div style={{ marginBottom: 56, paddingTop: 16 }}>
        <div className="fade-up" style={{ opacity: 0 }}>
          <p style={{
            fontFamily: 'var(--font-display)',
            fontSize: 13,
            letterSpacing: '0.2em',
            textTransform: 'uppercase',
            color: 'var(--gold)',
            marginBottom: 12,
          }}>
            Now Showing
          </p>
          <h1 className="page-title" style={{ marginBottom: 0 }}>
            What will you<br />
            <em style={{ fontStyle: 'italic', color: 'var(--silver)' }}>watch tonight?</em>
          </h1>
        </div>
      </div>

      {/* Demo notice */}
      {error && (
        <div style={{
          marginBottom: 24,
          padding: '10px 16px',
          background: 'rgba(201,168,76,0.08)',
          border: '1px solid var(--gold-dim)',
          borderRadius: 'var(--radius-sm)',
          fontSize: 13,
          color: 'var(--gold)',
        }}>
          ⚠ {error}
        </div>
      )}

      {/* Genre filter */}
      <div style={{
        display: 'flex',
        gap: 6,
        flexWrap: 'wrap',
        marginBottom: 40,
      }}>
        {GENRES.map(g => (
          <button
            key={g}
            onClick={() => setGenre(g)}
            style={{
              padding: '7px 16px',
              borderRadius: 20,
              border: `1px solid ${genre === g ? 'var(--gold)' : 'var(--ink-4)'}`,
              background: genre === g ? 'rgba(201,168,76,0.12)' : 'transparent',
              color: genre === g ? 'var(--gold)' : 'var(--silver)',
              fontSize: 13,
              cursor: 'pointer',
              transition: 'all 0.15s ease',
              letterSpacing: '0.02em',
            }}
          >
            {g}
          </button>
        ))}
      </div>

      {/* Grid */}
      {loading ? (
        <div className="spinner" />
      ) : filtered.length === 0 ? (
        <div style={{ textAlign: 'center', padding: '60px 0', color: 'var(--silver)' }}>
          <div style={{ fontFamily: 'var(--font-display)', fontSize: 32, marginBottom: 8 }}>No films</div>
          <p>No films found for this genre.</p>
        </div>
      ) : (
        <div className="grid-4">
          {filtered.map((movie, i) => (
            <MovieCard key={movie.id} movie={movie} delay={i * 50} />
          ))}
        </div>
      )}
    </div>
  );
}
