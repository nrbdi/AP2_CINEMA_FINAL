package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cinema/movie-service/internal/domain"
	"cinema/movie-service/internal/usecase"
	"github.com/google/uuid"
)

// ── Mocks ────────────────────────────────────────────────
type mockMovieRepo struct {
	movies map[uuid.UUID]*domain.Movie
}

func newMockMovieRepo() *mockMovieRepo {
	return &mockMovieRepo{movies: make(map[uuid.UUID]*domain.Movie)}
}
func (m *mockMovieRepo) Create(_ context.Context, mov *domain.Movie) error {
	m.movies[mov.ID] = mov
	return nil
}
func (m *mockMovieRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Movie, error) {
	mov, ok := m.movies[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return mov, nil
}
func (m *mockMovieRepo) List(_ context.Context, genre string, page, pageSize int) ([]*domain.Movie, int, error) {
	var result []*domain.Movie
	for _, mov := range m.movies {
		if genre == "" || mov.Genre == genre {
			result = append(result, mov)
		}
	}
	return result, len(result), nil
}
func (m *mockMovieRepo) Update(_ context.Context, mov *domain.Movie) error {
	if _, ok := m.movies[mov.ID]; !ok {
		return domain.ErrNotFound
	}
	m.movies[mov.ID] = mov
	return nil
}
func (m *mockMovieRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.movies, id)
	return nil
}

type mockHallRepo struct{}
func (m *mockHallRepo) Create(_ context.Context, h *domain.Hall, rows, spr int) (*domain.Hall, error) {
	h.ID = uuid.New()
	h.TotalSeats = rows * spr
	return h, nil
}
func (m *mockHallRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Hall, error) { return nil, nil }
func (m *mockHallRepo) GetSeatMap(_ context.Context, hallID uuid.UUID) ([]*domain.Seat, error) { return nil, nil }

type mockShowtimeRepo struct {
	showtimes map[uuid.UUID]*domain.Showtime
}
func newMockShowtimeRepo() *mockShowtimeRepo {
	return &mockShowtimeRepo{showtimes: make(map[uuid.UUID]*domain.Showtime)}
}
func (m *mockShowtimeRepo) Create(_ context.Context, s *domain.Showtime) error {
	m.showtimes[s.ID] = s; return nil
}
func (m *mockShowtimeRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Showtime, error) {
	s, ok := m.showtimes[id]
	if !ok { return nil, domain.ErrNotFound }
	return s, nil
}
func (m *mockShowtimeRepo) List(_ context.Context, movieID uuid.UUID, date string) ([]*domain.Showtime, error) {
	return nil, nil
}
func (m *mockShowtimeRepo) GetAvailableSeats(_ context.Context, showtimeID uuid.UUID) ([]*domain.Seat, error) {
	return nil, nil
}

type mockBookingRepo struct {
	bookings map[uuid.UUID]*domain.Booking
}
func newMockBookingRepo() *mockBookingRepo {
	return &mockBookingRepo{bookings: make(map[uuid.UUID]*domain.Booking)}
}
func (m *mockBookingRepo) CreateWithSeats(_ context.Context, b *domain.Booking) error {
	m.bookings[b.ID] = b; return nil
}
func (m *mockBookingRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Booking, error) {
	b, ok := m.bookings[id]
	if !ok { return nil, domain.ErrNotFound }
	return b, nil
}
func (m *mockBookingRepo) ListByUserID(_ context.Context, userID uuid.UUID) ([]*domain.Booking, error) {
	var result []*domain.Booking
	for _, b := range m.bookings {
		if b.UserID == userID { result = append(result, b) }
	}
	return result, nil
}
func (m *mockBookingRepo) Cancel(_ context.Context, bookingID, userID uuid.UUID) (*domain.Booking, error) {
	b, ok := m.bookings[bookingID]
	if !ok { return nil, domain.ErrNotFound }
	b.Status = domain.BookingStatusCancelled
	return b, nil
}

type mockPublisher struct{ published []domain.BookingConfirmedEvent }
func (p *mockPublisher) PublishBookingConfirmed(_ context.Context, e *domain.BookingConfirmedEvent) error {
	p.published = append(p.published, *e); return nil
}
func (p *mockPublisher) PublishBookingCancelled(_ context.Context, e *domain.BookingCancelledEvent) error { return nil }

type mockCache struct{}
func (c *mockCache) SetAvailableSeats(_ context.Context, _ string, _ []*domain.Seat) error { return nil }
func (c *mockCache) GetAvailableSeats(_ context.Context, _ string) ([]*domain.Seat, bool, error) { return nil, false, nil }
func (c *mockCache) InvalidateShowtime(_ context.Context, _ string) error { return nil }

func newTestUsecase() (*usecase.MovieUsecase, *mockMovieRepo, *mockBookingRepo, *mockPublisher, *mockShowtimeRepo) {
	movieRepo    := newMockMovieRepo()
	bookingRepo  := newMockBookingRepo()
	pub          := &mockPublisher{}
	showtimeRepo := newMockShowtimeRepo()
	uc := usecase.NewMovieUsecase(movieRepo, &mockHallRepo{}, showtimeRepo, bookingRepo, pub, &mockCache{})
	return uc, movieRepo, bookingRepo, pub, showtimeRepo
}

// ── Tests ────────────────────────────────────────────────
func TestCreateMovie_Success(t *testing.T) {
	uc, _, _, _, _ := newTestUsecase()
	m, err := uc.CreateMovie(context.Background(), &domain.Movie{
		Title: "Inception", Genre: "Sci-Fi", DurationMin: 148,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if m.ID == uuid.Nil {
		t.Error("expected non-nil UUID")
	}
	if m.Title != "Inception" {
		t.Errorf("expected title 'Inception', got %q", m.Title)
	}
}

func TestCreateMovie_EmptyTitle(t *testing.T) {
	uc, _, _, _, _ := newTestUsecase()
	_, err := uc.CreateMovie(context.Background(), &domain.Movie{})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestGetMovie_NotFound(t *testing.T) {
	uc, _, _, _, _ := newTestUsecase()
	_, err := uc.GetMovie(context.Background(), uuid.New())
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCreateBooking_Success(t *testing.T) {
	uc, _, bookingRepo, pub, showtimeRepo := newTestUsecase()

	showtime := &domain.Showtime{
		ID: uuid.New(), MovieID: uuid.New(), HallID: uuid.New(),
		StartTime: time.Now().Add(2 * time.Hour),
		EndTime:   time.Now().Add(4 * time.Hour),
		Price:     12.50,
	}
	showtimeRepo.showtimes[showtime.ID] = showtime

	userID  := uuid.New()
	seatIDs := []uuid.UUID{uuid.New(), uuid.New()}

	booking, err := uc.CreateBooking(context.Background(), userID, showtime.ID, seatIDs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if booking.TotalPrice != 25.0 {
		t.Errorf("expected total_price 25.0, got %.2f", booking.TotalPrice)
	}
	if len(pub.published) != 1 {
		t.Errorf("expected 1 published event, got %d", len(pub.published))
	}
	if _, ok := bookingRepo.bookings[booking.ID]; !ok {
		t.Error("booking not found in repo")
	}
}

func TestCreateBooking_NoSeats(t *testing.T) {
	uc, _, _, _, _ := newTestUsecase()
	_, err := uc.CreateBooking(context.Background(), uuid.New(), uuid.New(), nil)
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestListMovies_Pagination(t *testing.T) {
	uc, movieRepo, _, _, _ := newTestUsecase()

	for i := 0; i < 5; i++ {
		id := uuid.New()
		movieRepo.movies[id] = &domain.Movie{ID: id, Title: "Movie", Genre: "Action"}
	}

	movies, total, err := uc.ListMovies(context.Background(), "Action", 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(movies) != 5 {
		t.Errorf("expected 5 movies, got %d", len(movies))
	}
}
