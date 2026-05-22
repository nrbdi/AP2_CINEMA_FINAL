import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { api } from '../services/api';

export default function ProfilePage() {
  const { user, logout } = useAuth();
  const [form, setForm] = useState({ full_name: user?.full_name || '', phone: user?.phone || '' });
  const [saving, setSaving] = useState(false);
  const [msg, setMsg] = useState('');

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    setMsg('');
    try {
      await api.updateProfile(form);
      setMsg('Profile updated.');
    } catch (err) {
      setMsg(err.message || 'Could not update profile.');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="page" style={{ maxWidth: 600 }}>
      <h1 className="page-title fade-up" style={{ opacity: 0 }}>Profile</h1>
      <p className="page-sub fade-up fade-up-1">Manage your account details.</p>

      {/* Account card */}
      <div
        className="fade-up fade-up-2"
        style={{
          opacity: 0,
          background: 'var(--ink-2)',
          border: '1px solid var(--ink-4)',
          borderRadius: 'var(--radius-lg)',
          padding: '28px 32px',
          marginBottom: 20,
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: 16, marginBottom: 28 }}>
          <div style={{
            width: 52, height: 52,
            borderRadius: '50%',
            background: 'rgba(201,168,76,0.15)',
            border: '1px solid var(--gold-dim)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontFamily: 'var(--font-display)',
            fontSize: 22,
            color: 'var(--gold)',
            fontWeight: 300,
          }}>
            {user?.full_name?.charAt(0) || '?'}
          </div>
          <div>
            <div style={{ fontFamily: 'var(--font-display)', fontSize: 22, fontWeight: 300 }}>
              {user?.full_name}
            </div>
            <div style={{ fontSize: 13, color: 'var(--silver)' }}>{user?.email}</div>
          </div>
          <span className={`badge badge-${user?.role === 'admin' ? 'gold' : 'silver'}`} style={{ marginLeft: 'auto' }}>
            {user?.role || 'user'}
          </span>
        </div>

        <form onSubmit={handleSave} style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
          <div className="grid-2">
            <div>
              <label style={{ display: 'block', fontSize: 12, color: 'var(--silver)', marginBottom: 6, letterSpacing: '0.06em', textTransform: 'uppercase' }}>
                Full name
              </label>
              <input
                className="input"
                value={form.full_name}
                onChange={e => setForm({ ...form, full_name: e.target.value })}
              />
            </div>
            <div>
              <label style={{ display: 'block', fontSize: 12, color: 'var(--silver)', marginBottom: 6, letterSpacing: '0.06em', textTransform: 'uppercase' }}>
                Phone
              </label>
              <input
                className="input"
                value={form.phone}
                onChange={e => setForm({ ...form, phone: e.target.value })}
                placeholder="optional"
              />
            </div>
          </div>

          <div>
            <label style={{ display: 'block', fontSize: 12, color: 'var(--silver)', marginBottom: 6, letterSpacing: '0.06em', textTransform: 'uppercase' }}>
              Email
            </label>
            <input className="input" value={user?.email || ''} disabled style={{ opacity: 0.5 }} />
          </div>

          {msg && (
            <div style={{
              padding: '10px 14px',
              background: msg.includes('updated') ? 'rgba(75,175,122,0.1)' : 'rgba(200,75,75,0.1)',
              border: `1px solid ${msg.includes('updated') ? 'rgba(75,175,122,0.3)' : 'rgba(200,75,75,0.3)'}`,
              borderRadius: 'var(--radius-sm)',
              fontSize: 13,
              color: msg.includes('updated') ? 'var(--green)' : 'var(--red)',
            }}>
              {msg}
            </div>
          )}

          <div style={{ display: 'flex', gap: 12 }}>
            <button type="submit" className="btn btn-gold" disabled={saving}>
              {saving ? 'Saving…' : 'Save changes'}
            </button>
          </div>
        </form>
      </div>

      {/* Danger zone */}
      <div
        className="fade-up fade-up-3"
        style={{
          opacity: 0,
          background: 'var(--ink-2)',
          border: '1px solid rgba(200,75,75,0.15)',
          borderRadius: 'var(--radius-lg)',
          padding: '20px 28px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          gap: 16,
        }}
      >
        <div>
          <div style={{ fontSize: 14, marginBottom: 2 }}>Sign out of your account</div>
          <div style={{ fontSize: 12, color: 'var(--silver)' }}>You'll need to sign in again to access bookings.</div>
        </div>
        <button className="btn btn-danger" onClick={logout} style={{ flexShrink: 0 }}>
          Sign out
        </button>
      </div>
    </div>
  );
}
