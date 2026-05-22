import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [form, setForm] = useState({ email: '', password: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await login(form.email, form.password);
      navigate('/');
    } catch (err) {
      setError(err.message || 'Invalid email or password.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{
      minHeight: 'calc(100vh - 64px)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '40px 24px',
    }}>
      <div style={{ width: '100%', maxWidth: 400 }}>
        <div className="fade-up" style={{ opacity: 0, textAlign: 'center', marginBottom: 40 }}>
          <h1 style={{
            fontFamily: 'var(--font-display)',
            fontSize: 42,
            fontWeight: 300,
            marginBottom: 8,
          }}>
            Welcome back
          </h1>
          <p style={{ color: 'var(--silver)', fontSize: 14 }}>
            Sign in to access your bookings
          </p>
        </div>

        <div
          className="fade-up fade-up-1"
          style={{
            opacity: 0,
            background: 'var(--ink-2)',
            border: '1px solid var(--ink-4)',
            borderRadius: 'var(--radius-lg)',
            padding: '32px',
          }}
        >
          <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
            <div>
              <label style={{ display: 'block', fontSize: 12, color: 'var(--silver)', marginBottom: 6, letterSpacing: '0.06em', textTransform: 'uppercase' }}>
                Email
              </label>
              <input
                className="input"
                type="email"
                required
                value={form.email}
                onChange={e => setForm({ ...form, email: e.target.value })}
                placeholder="you@example.com"
              />
            </div>

            <div>
              <label style={{ display: 'block', fontSize: 12, color: 'var(--silver)', marginBottom: 6, letterSpacing: '0.06em', textTransform: 'uppercase' }}>
                Password
              </label>
              <input
                className="input"
                type="password"
                required
                value={form.password}
                onChange={e => setForm({ ...form, password: e.target.value })}
                placeholder="••••••••"
              />
            </div>

            {error && (
              <div style={{
                padding: '10px 14px',
                background: 'rgba(200,75,75,0.1)',
                border: '1px solid rgba(200,75,75,0.3)',
                borderRadius: 'var(--radius-sm)',
                fontSize: 13,
                color: 'var(--red)',
              }}>
                {error}
              </div>
            )}

            <button
              type="submit"
              className="btn btn-gold"
              style={{ width: '100%', justifyContent: 'center', padding: '13px', marginTop: 4 }}
              disabled={loading}
            >
              {loading ? 'Signing in…' : 'Sign in'}
            </button>
          </form>

          <div style={{
            marginTop: 24,
            paddingTop: 24,
            borderTop: '1px solid var(--ink-4)',
            textAlign: 'center',
            fontSize: 13,
            color: 'var(--silver)',
          }}>
            Don't have an account?{' '}
            <Link to="/register" style={{ color: 'var(--gold)', textDecoration: 'none' }}>
              Join now
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
