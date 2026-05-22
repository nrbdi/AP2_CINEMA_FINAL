import { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function Navbar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [menuOpen, setMenuOpen] = useState(false);

  const handleLogout = async () => {
    await logout();
    navigate('/');
  };

  const isActive = (path) => location.pathname === path;

  return (
    <nav style={{
      background: 'rgba(10,10,11,0.92)',
      borderBottom: '1px solid #1A1A1F',
      backdropFilter: 'blur(20px)',
      position: 'sticky',
      top: 0,
      zIndex: 100,
    }}>
      <div style={{
        maxWidth: 1200,
        margin: '0 auto',
        padding: '0 24px',
        height: 64,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
      }}>
        {/* Logo */}
        <Link to="/" style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
          <span style={{
            fontFamily: 'var(--font-display)',
            fontSize: 22,
            fontWeight: 300,
            letterSpacing: '0.12em',
            color: 'var(--white)',
          }}>
            CINEMA
          </span>
          <span style={{
            width: 6, height: 6,
            borderRadius: '50%',
            background: 'var(--gold)',
            display: 'inline-block',
            marginBottom: 2,
          }} />
        </Link>

        {/* Desktop nav */}
        <div style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
          <NavLink to="/" active={isActive('/')}>Films</NavLink>

          {user ? (
            <>
              <NavLink to="/my-bookings" active={isActive('/my-bookings')}>My Bookings</NavLink>
              <NavLink to="/profile" active={isActive('/profile')}>{user.full_name || 'Profile'}</NavLink>
              <button
                onClick={handleLogout}
                style={{
                  marginLeft: 8,
                  padding: '7px 18px',
                  background: 'transparent',
                  border: '1px solid var(--smoke)',
                  borderRadius: 'var(--radius-sm)',
                  color: 'var(--silver)',
                  fontSize: 12,
                  letterSpacing: '0.06em',
                  textTransform: 'uppercase',
                  cursor: 'pointer',
                  transition: 'var(--transition)',
                }}
                onMouseEnter={e => { e.target.style.borderColor = 'var(--mist)'; e.target.style.color = 'var(--white)'; }}
                onMouseLeave={e => { e.target.style.borderColor = 'var(--smoke)'; e.target.style.color = 'var(--silver)'; }}
              >
                Sign out
              </button>
            </>
          ) : (
            <>
              <NavLink to="/login" active={isActive('/login')}>Sign in</NavLink>
              <Link to="/register">
                <button className="btn btn-gold" style={{ marginLeft: 8, padding: '8px 20px' }}>
                  Join
                </button>
              </Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
}

function NavLink({ to, children, active }) {
  return (
    <Link to={to} style={{
      padding: '8px 14px',
      borderRadius: 'var(--radius-sm)',
      fontSize: 13,
      letterSpacing: '0.04em',
      color: active ? 'var(--gold)' : 'var(--silver)',
      background: active ? 'rgba(201,168,76,0.08)' : 'transparent',
      transition: 'var(--transition)',
      textDecoration: 'none',
    }}
    onMouseEnter={e => { if (!active) e.currentTarget.style.color = 'var(--white)'; }}
    onMouseLeave={e => { if (!active) e.currentTarget.style.color = 'var(--silver)'; }}
    >
      {children}
    </Link>
  );
}
