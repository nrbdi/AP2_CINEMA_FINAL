package usecase

import (
	"context"
	"fmt"
	"time"

	"cinema/movie-service/internal/domain"
	"github.com/google/uuid"
)

type MovieUsecase struct {
	movieRepo    domain.MovieRepository
	hallRepo     domain.HallRepository
	showtimeRepo domain.ShowtimeRepository
	bookingRepo  domain.BookingRepository
	publisher    domain.Publisher
	cache        domain.SeatCache
}

func NewMovieUsecase(
	mr domain.MovieRepository,
	hr domain.HallRepository,
	sr domain.ShowtimeRepository,
	br domain.BookingRepository,
	pub domain.Publisher,
	cache domain.SeatCache,
) *MovieUsecase {
	return &MovieUsecase{
		movieRepo:    mr,
		hallRepo:     hr,
		showtimeRepo: sr,
		bookingRepo:  br,
		publisher:    pub,
		cache:        cache,
	}
}

// ── Movie ────────────────────────────────────────────────
func (uc *MovieUsecase) CreateMovie(ctx context.Context, m *domain.Movie) (*domain.Movie, error) {
	if m.Title == "" {
		return nil, fmt.Errorf("%w: title is required", domain.ErrInvalidInput)
	}
	m.ID = uuid.New()
	m.CreatedAt = time.Now()
	if err := uc.movieRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (uc *MovieUsecase) GetMovie(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	return uc.movieRepo.GetByID(ctx, id)
}

func (uc *MovieUsecase) ListMovies(ctx context.Context, genre string, page, pageSize int) ([]*domain.Movie, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return uc.movieRepo.List(ctx, genre, page, pageSize)
}

func (uc *MovieUsecase) UpdateMovie(ctx context.Context, m *domain.Movie) (*domain.Movie, error) {
	existing, err := uc.movieRepo.GetByID(ctx, m.ID)
	if err != nil {
		return nil, err
	}
	if m.Title != "" {
		existing.Title = m.Title
	}
	if m.Description != "" {
		existing.Description = m.Description
	}
	if m.DurationMin > 0 {
		existing.DurationMin = m.DurationMin
	}
	if m.Genre != "" {
		existing.Genre = m.Genre
	}
	if m.PosterURL != "" {
		existing.PosterURL = m.PosterURL
	}
	if err := uc.movieRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (uc *MovieUsecase) DeleteMovie(ctx context.Context, id uuid.UUID) error {
	return uc.movieRepo.Delete(ctx, id)
}

// ── Hall ─────────────────────────────────────────────────
func (uc *MovieUsecase) CreateHall(ctx context.Context, name string, rows, seatsPerRow int) (*domain.Hall, error) {
	if name == "" || rows < 1 || seatsPerRow < 1 {
		return nil, fmt.Errorf("%w: invalid hall parameters", domain.ErrInvalidInput)
	}
	h := &domain.Hall{Name: name}
	return uc.hallRepo.Create(ctx, h, rows, seatsPerRow)
}

func (uc *MovieUsecase) GetHallSeatMap(ctx context.Context, hallID uuid.UUID) ([]*domain.Seat, error) {
	return uc.hallRepo.GetSeatMap(ctx, hallID)
}

// ── Showtime ─────────────────────────────────────────────
func (uc *MovieUsecase) CreateShowtime(ctx context.Context, s *domain.Showtime) (*domain.Showtime, error) {
	s.ID = uuid.New()
	if err := uc.showtimeRepo.Create(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}

func (uc *MovieUsecase) ListShowtimes(ctx context.Context, movieID uuid.UUID, date string) ([]*domain.Showtime, error) {
	return uc.showtimeRepo.List(ctx, movieID, date)
}

func (uc *MovieUsecase) GetAvailableSeats(ctx context.Context, showtimeID uuid.UUID) ([]*domain.Seat, error) {
	// Try cache first
	seats, hit, err := uc.cache.GetAvailableSeats(ctx, showtimeID.String())
	if err == nil && hit {
		return seats, nil
	}
	// Cache miss — query DB
	seats, err = uc.showtimeRepo.GetAvailableSeats(ctx, showtimeID)
	if err != nil {
		return nil, err
	}
	// Populate cache (best-effort)
	_ = uc.cache.SetAvailableSeats(ctx, showtimeID.String(), seats)
	return seats, nil
}

// ── Booking ──────────────────────────────────────────────
func (uc *MovieUsecase) CreateBooking(ctx context.Context, userID, showtimeID uuid.UUID, seatIDs []uuid.UUID) (*domain.Booking, error) {
	if len(seatIDs) == 0 {
		return nil, fmt.Errorf("%w: at least one seat required", domain.ErrInvalidInput)
	}

	showtime, err := uc.showtimeRepo.GetByID(ctx, showtimeID)
	if err != nil {
		return nil, err
	}

	totalPrice := showtime.Price * float64(len(seatIDs))

	booking := &domain.Booking{
		ID:         uuid.New(),
		UserID:     userID,
		ShowtimeID: showtimeID,
		Status:     domain.BookingStatusConfirmed,
		SeatIDs:    seatIDs,
		TotalPrice: totalPrice,
		BookedAt:   time.Now(),
	}

	// This runs in a DB transaction with FOR UPDATE locking
	if err := uc.bookingRepo.CreateWithSeats(ctx, booking); err != nil {
		return nil, err
	}

	// Invalidate seat cache for this showtime
	_ = uc.cache.InvalidateShowtime(ctx, showtimeID.String())

	// Publish event for payment-service
	seatIDStrs := make([]string, len(seatIDs))
	for i, id := range seatIDs {
		seatIDStrs[i] = id.String()
	}
	_ = uc.publisher.PublishBookingConfirmed(ctx, &domain.BookingConfirmedEvent{
		BookingID:  booking.ID.String(),
		UserID:     userID.String(),
		ShowtimeID: showtimeID.String(),
		SeatIDs:    seatIDStrs,
		TotalPrice: totalPrice,
	})

	return booking, nil
}

func (uc *MovieUsecase) GetBooking(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error) {
	return uc.bookingRepo.GetByID(ctx, bookingID)
}

func (uc *MovieUsecase) ListUserBookings(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error) {
	return uc.bookingRepo.ListByUserID(ctx, userID)
}

func (uc *MovieUsecase) CancelBooking(ctx context.Context, bookingID, userID uuid.UUID) (*domain.Booking, error) {
	booking, err := uc.bookingRepo.Cancel(ctx, bookingID, userID)
	if err != nil {
		return nil, err
	}
	_ = uc.cache.InvalidateShowtime(ctx, booking.ShowtimeID.String())
	_ = uc.publisher.PublishBookingCancelled(ctx, &domain.BookingCancelledEvent{
		BookingID: bookingID.String(),
		UserID:    userID.String(),
	})
	return booking, nil
}
