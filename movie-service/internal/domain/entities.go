package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ── Errors ───────────────────────────────────────────────
var (
	ErrNotFound          = errors.New("not found")
	ErrSeatsUnavailable  = errors.New("one or more seats are unavailable")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidInput      = errors.New("invalid input")
)

// ── Movie ────────────────────────────────────────────────
type Movie struct {
	ID          uuid.UUID
	Title       string
	Description string
	DurationMin int
	Genre       string
	PosterURL   string
	ReleaseDate string
	CreatedAt   time.Time
}

// ── Hall ─────────────────────────────────────────────────
type Hall struct {
	ID         uuid.UUID
	Name       string
	TotalSeats int
	CreatedAt  time.Time
}

// ── Seat ─────────────────────────────────────────────────
type SeatType string

const (
	SeatTypeStandard SeatType = "standard"
	SeatTypePremium  SeatType = "premium"
	SeatTypeVIP      SeatType = "vip"
)

type Seat struct {
	ID         uuid.UUID
	HallID     uuid.UUID
	RowNumber  int
	SeatNumber int
	SeatType   SeatType
}

// ── Showtime ─────────────────────────────────────────────
type Showtime struct {
	ID        uuid.UUID
	MovieID   uuid.UUID
	HallID    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	Price     float64
}

// ── Booking ──────────────────────────────────────────────
type BookingStatus string

const (
	BookingStatusConfirmed  BookingStatus = "confirmed"
	BookingStatusCancelled  BookingStatus = "cancelled"
)

type Booking struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	ShowtimeID  uuid.UUID
	Status      BookingStatus
	SeatIDs     []uuid.UUID
	TotalPrice  float64
	BookedAt    time.Time
	CancelledAt *time.Time
}

// ── NATS Events ──────────────────────────────────────────
type BookingConfirmedEvent struct {
	BookingID  string   `json:"booking_id"`
	UserID     string   `json:"user_id"`
	ShowtimeID string   `json:"showtime_id"`
	SeatIDs    []string `json:"seat_ids"`
	TotalPrice float64  `json:"total_price"`
}

type BookingCancelledEvent struct {
	BookingID string `json:"booking_id"`
	UserID    string `json:"user_id"`
}
