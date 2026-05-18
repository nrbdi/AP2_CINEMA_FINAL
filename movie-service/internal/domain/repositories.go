package domain

import (
	"context"

	"github.com/google/uuid"
)

type MovieRepository interface {
	Create(ctx context.Context, m *Movie) error
	GetByID(ctx context.Context, id uuid.UUID) (*Movie, error)
	List(ctx context.Context, genre string, page, pageSize int) ([]*Movie, int, error)
	Update(ctx context.Context, m *Movie) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type HallRepository interface {
	Create(ctx context.Context, h *Hall, rows, seatsPerRow int) (*Hall, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Hall, error)
	GetSeatMap(ctx context.Context, hallID uuid.UUID) ([]*Seat, error)
}

type ShowtimeRepository interface {
	Create(ctx context.Context, s *Showtime) error
	GetByID(ctx context.Context, id uuid.UUID) (*Showtime, error)
	List(ctx context.Context, movieID uuid.UUID, date string) ([]*Showtime, error)
	GetAvailableSeats(ctx context.Context, showtimeID uuid.UUID) ([]*Seat, error)
}

type BookingRepository interface {
	CreateWithSeats(ctx context.Context, b *Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*Booking, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Booking, error)
	Cancel(ctx context.Context, bookingID, userID uuid.UUID) (*Booking, error)
}

type Publisher interface {
	PublishBookingConfirmed(ctx context.Context, event *BookingConfirmedEvent) error
	PublishBookingCancelled(ctx context.Context, event *BookingCancelledEvent) error
}

type SeatCache interface {
	SetAvailableSeats(ctx context.Context, showtimeID string, seats []*Seat) error
	GetAvailableSeats(ctx context.Context, showtimeID string) ([]*Seat, bool, error)
	InvalidateShowtime(ctx context.Context, showtimeID string) error
}
