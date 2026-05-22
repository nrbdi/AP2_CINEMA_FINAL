import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function RegisterPage() {
  const { register, login } = useAuth();
  const navigate = useNavigate();
  const [form, setForm] = useState({ email: '', password: '', full_name: '', phone: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const set = (k) => (e) => setForm({ ...form, [k]: e.target.value });

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    if (form.password.length < 8) { setError('Password must be at least 8 characters.'); return; }
    setLoading(true);
    try {
      await register(form);
      await login(form.email, form.password);
      navigate('/');
    } catch (err) {
      setError(err.message || 'Registration failed.');
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
      <div style={{ width: '100%', maxWidth: 440 }}>
        <div className="fade-up" style={{ opacity: 0, textAlign: 'center', marginBottom: 40 }}>
          <h1 style={{
            fontFamily: 'var(--font-display)',
            fontSize: 42,
            fontWeight: 300,
            marginBottom: 8,
          }}>
            Create account
          </h1>
          <p style={{ color: 'var(--silver)', fontSize: 14 }}>
            Join to start booking cinema tickets
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
            {[
              { key: 'full_name', label: 'Full name', type: 'text', placeholder: 'Alice Smith', required: true },
              { key: 'email', label: 'Email', type: 'email', placeholder: 'alice@example.com', required: true },
              { key: 'phone', label: 'Phone', type: 'tel', placeholder: '+1 555 0100', required: false },
              { key: 'password', label: 'Password', type: 'password', placeholder: 'min 8 characters', required: true },
            ].map(({ key, label, type, placeholder, required }) => (
              <div key={key}>
                <label style={{
                  display: 'block',
                  fontSize: 12,
                  color: 'var(--silver)',
                  marginBottom: 6,
                  letterSpacing: '0.06em',
                  textTransform: 'uppercase',
                }}>
                  {label}{!required && <span style={{ color: 'var(--mist)', marginLeft: 4 }}>(optional)</span>}
                </label>
                <input
                  className="input"
                  type={type}
                  required={required}
                  placeholder={placeholder}
                  value={form[key]}
                  onChange={set(key)}
                />
              </div>
            ))}

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
              {loading ? 'Creating account…' : 'Create account'}
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
            Already have an account?{' '}
            <Link to="/login" style={{ color: 'var(--gold)', textDecoration: 'none' }}>
              Sign in
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
