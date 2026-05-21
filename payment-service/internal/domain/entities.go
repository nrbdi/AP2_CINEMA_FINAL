package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type Payment struct {
	ID        uuid.UUID
	BookingID uuid.UUID
	UserID    uuid.UUID
	Amount    float64
	Status    PaymentStatus
	Method    string
	PaidAt    time.Time
}

type BookingConfirmedEvent struct {
	BookingID  string   `json:"booking_id"`
	UserID     string   `json:"user_id"`
	ShowtimeID string   `json:"showtime_id"`
	SeatIDs    []string `json:"seat_ids"`
	TotalPrice float64  `json:"total_price"`
}

type PaymentRepository interface {
	Create(ctx context.Context, p *Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*Payment, error)
	GetByBookingID(ctx context.Context, bookingID uuid.UUID) (*Payment, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Payment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status PaymentStatus) error
}

type EmailSender interface {
	SendPaymentReceipt(toEmail, toName, movieTitle, showtime string, seats []string, amount float64) error
}
