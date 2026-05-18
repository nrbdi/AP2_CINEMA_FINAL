package postgres

import (
	"context"

	"cinema/payment-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepo struct{ db *pgxpool.Pool }

func NewPaymentRepo(db *pgxpool.Pool) *PaymentRepo { return &PaymentRepo{db: db} }

func (r *PaymentRepo) Create(ctx context.Context, p *domain.Payment) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO payments (id,booking_id,user_id,amount,status,method,paid_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		p.ID, p.BookingID, p.UserID, p.Amount, p.Status, p.Method, p.PaidAt,
	)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PaymentRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	p := &domain.Payment{}
	err := r.db.QueryRow(ctx,
		`SELECT id,booking_id,user_id,amount,status,method,paid_at FROM payments WHERE id=$1`, id,
	).Scan(&p.ID, &p.BookingID, &p.UserID, &p.Amount, &p.Status, &p.Method, &p.PaidAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return p, err
}

func (r *PaymentRepo) GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*domain.Payment, error) {
	p := &domain.Payment{}
	err := r.db.QueryRow(ctx,
		`SELECT id,booking_id,user_id,amount,status,method,paid_at FROM payments WHERE booking_id=$1`, bookingID,
	).Scan(&p.ID, &p.BookingID, &p.UserID, &p.Amount, &p.Status, &p.Method, &p.PaidAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return p, err
}

func (r *PaymentRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Payment, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,booking_id,user_id,amount,status,method,paid_at FROM payments WHERE user_id=$1 ORDER BY paid_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []*domain.Payment
	for rows.Next() {
		p := &domain.Payment{}
		if err := rows.Scan(&p.ID, &p.BookingID, &p.UserID, &p.Amount, &p.Status, &p.Method, &p.PaidAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *PaymentRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus) error {
	_, err := r.db.Exec(ctx, `UPDATE payments SET status=$1 WHERE id=$2`, status, id)
	return err
}
