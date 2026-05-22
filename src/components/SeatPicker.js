import { useState } from 'react';

export default function SeatPicker({ seats, onSelectionChange }) {
  const [selected, setSelected] = useState(new Set());

  if (!seats || seats.length === 0) {
    return (
      <div style={{ textAlign: 'center', padding: '40px', color: 'var(--silver)' }}>
        No seats available for this showtime.
      </div>
    );
  }

  // Group by row
  const rows = {};
  seats.forEach(s => {
    if (!rows[s.row_number]) rows[s.row_number] = [];
    rows[s.row_number].push(s);
  });

  const toggle = (seatId) => {
    const next = new Set(selected);
    if (next.has(seatId)) next.delete(seatId);
    else next.add(seatId);
    setSelected(next);
    onSelectionChange([...next]);
  };

  const typeColor = { standard: 'var(--smoke)', premium: '#185FA542', vip: '#C9A84C22' };
  const typeBorder = { standard: 'var(--mist)', premium: '#185FA5', vip: 'var(--gold-dim)' };

  return (
    <div>
      {/* Screen indicator */}
      <div style={{ textAlign: 'center', marginBottom: 32 }}>
        <div style={{
          display: 'inline-block',
          padding: '6px 60px',
          background: 'linear-gradient(180deg, var(--smoke) 0%, transparent 100%)',
          borderRadius: '0 0 50% 50% / 0 0 16px 16px',
          fontSize: 11,
          color: 'var(--silver)',
          letterSpacing: '0.14em',
          textTransform: 'uppercase',
          marginBottom: 8,
        }}>
          Screen
        </div>
      </div>

      {/* Seat grid */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: 8, alignItems: 'center' }}>
        {Object.entries(rows).sort(([a],[b]) => a-b).map(([rowNum, rowSeats]) => (
          <div key={rowNum} style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
            <span style={{
              width: 20, fontSize: 11, color: 'var(--mist)',
              textAlign: 'right', flexShrink: 0,
            }}>
              {rowNum}
            </span>
            <div style={{ display: 'flex', gap: 5 }}>
              {rowSeats.sort((a,b) => a.seat_number - b.seat_number).map(seat => {
                const isSel = selected.has(seat.id);
                return (
                  <button
                    key={seat.id}
                    onClick={() => toggle(seat.id)}
                    title={`Row ${seat.row_number}, Seat ${seat.seat_number} (${seat.seat_type})`}
                    style={{
                      width: 28, height: 26,
                      borderRadius: '4px 4px 2px 2px',
                      border: `1px solid ${isSel ? 'var(--gold)' : typeBorder[seat.seat_type] || 'var(--mist)'}`,
                      background: isSel ? 'var(--gold)' : (typeColor[seat.seat_type] || 'var(--smoke)'),
                      cursor: 'pointer',
                      transition: 'all 0.15s ease',
                      fontSize: 9,
                      color: isSel ? 'var(--ink)' : 'transparent',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                    }}
                    onMouseEnter={e => {
                      if (!isSel) {
                        e.currentTarget.style.background = 'var(--mist)';
                        e.currentTarget.style.borderColor = 'var(--silver)';
                      }
                    }}
                    onMouseLeave={e => {
                      if (!isSel) {
                        e.currentTarget.style.background = typeColor[seat.seat_type] || 'var(--smoke)';
                        e.currentTarget.style.borderColor = typeBorder[seat.seat_type] || 'var(--mist)';
                      }
                    }}
                  >
                    {seat.seat_number}
                  </button>
                );
              })}
            </div>
          </div>
        ))}
      </div>

      {/* Legend */}
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        gap: 24,
        marginTop: 28,
        fontSize: 12,
        color: 'var(--silver)',
      }}>
        {[
          { label: 'Standard', color: 'var(--smoke)', border: 'var(--mist)' },
          { label: 'Premium', color: '#185FA542', border: '#185FA5' },
          { label: 'VIP', color: '#C9A84C22', border: 'var(--gold-dim)' },
          { label: 'Selected', color: 'var(--gold)', border: 'var(--gold)' },
        ].map(({ label, color, border }) => (
          <div key={label} style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
            <div style={{
              width: 16, height: 14,
              borderRadius: '3px 3px 1px 1px',
              background: color,
              border: `1px solid ${border}`,
            }} />
            {label}
          </div>
        ))}
      </div>

      {/* Selection summary */}
      {selected.size > 0 && (
        <div style={{
          marginTop: 20,
          padding: '12px 16px',
          background: 'rgba(201,168,76,0.08)',
          border: '1px solid var(--gold-dim)',
          borderRadius: 'var(--radius-sm)',
          fontSize: 13,
          color: 'var(--gold)',
          textAlign: 'center',
        }}>
          {selected.size} seat{selected.size > 1 ? 's' : ''} selected
        </div>
      )}
    </div>
  );
}
