package postgres

import (
	"context"
	"fmt"
	"time"

	"cinema/movie-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ── Movie Repository ─────────────────────────────────────
type MovieRepo struct{ db *pgxpool.Pool }

func NewMovieRepo(db *pgxpool.Pool) *MovieRepo { return &MovieRepo{db: db} }

func (r *MovieRepo) Create(ctx context.Context, m *domain.Movie) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO movies (id,title,description,duration_min,genre,poster_url,release_date,created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		m.ID, m.Title, m.Description, m.DurationMin, m.Genre, m.PosterURL, m.ReleaseDate, m.CreatedAt,
	)
	return err
}

func (r *MovieRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	m := &domain.Movie{}
	var releaseDate *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT id,title,description,duration_min,genre,poster_url,release_date,created_at FROM movies WHERE id=$1`,
		id,
	).Scan(&m.ID, &m.Title, &m.Description, &m.DurationMin, &m.Genre, &m.PosterURL, &releaseDate, &m.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if releaseDate != nil {
		m.ReleaseDate = releaseDate.Format("2006-01-02")
	}
	return m, err
}

func (r *MovieRepo) List(ctx context.Context, genre string, page, pageSize int) ([]*domain.Movie, int, error) {
	offset := (page - 1) * pageSize
	var rows pgx.Rows
	var err error
	var total int

	if genre != "" {
		err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM movies WHERE genre=$1`, genre).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
		rows, err = r.db.Query(ctx,
			`SELECT id,title,description,duration_min,genre,poster_url,release_date,created_at
			 FROM movies WHERE genre=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
			genre, pageSize, offset,
		)
	} else {
		err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM movies`).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
		rows, err = r.db.Query(ctx,
			`SELECT id,title,description,duration_min,genre,poster_url,release_date,created_at
			 FROM movies ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
			pageSize, offset,
		)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []*domain.Movie
	for rows.Next() {
		m := &domain.Movie{}
		var rd *time.Time
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.DurationMin, &m.Genre, &m.PosterURL, &rd, &m.CreatedAt); err != nil {
			return nil, 0, err
		}
		if rd != nil {
			m.ReleaseDate = rd.Format("2006-01-02")
		}
		movies = append(movies, m)
	}
	return movies, total, nil
}

func (r *MovieRepo) Update(ctx context.Context, m *domain.Movie) error {
	_, err := r.db.Exec(ctx,
		`UPDATE movies SET title=$1,description=$2,duration_min=$3,genre=$4,poster_url=$5 WHERE id=$6`,
		m.Title, m.Description, m.DurationMin, m.Genre, m.PosterURL, m.ID,
	)
	return err
}

func (r *MovieRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM movies WHERE id=$1`, id)
	return err
}

// ── Hall Repository ──────────────────────────────────────
type HallRepo struct{ db *pgxpool.Pool }

func NewHallRepo(db *pgxpool.Pool) *HallRepo { return &HallRepo{db: db} }

func (r *HallRepo) Create(ctx context.Context, h *domain.Hall, rows, seatsPerRow int) (*domain.Hall, error) {
	h.ID = uuid.New()
	h.TotalSeats = rows * seatsPerRow
	h.CreatedAt = time.Now()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO halls (id,name,total_seats,created_at) VALUES ($1,$2,$3,$4)`,
		h.ID, h.Name, h.TotalSeats, h.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	for row := 1; row <= rows; row++ {
		for seat := 1; seat <= seatsPerRow; seat++ {
			_, err = tx.Exec(ctx,
				`INSERT INTO seats (id,hall_id,row_number,seat_number,seat_type) VALUES ($1,$2,$3,$4,$5)`,
				uuid.New(), h.ID, row, seat, domain.SeatTypeStandard,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	return h, tx.Commit(ctx)
}

func (r *HallRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Hall, error) {
	h := &domain.Hall{}
	err := r.db.QueryRow(ctx,
		`SELECT id,name,total_seats,created_at FROM halls WHERE id=$1`, id,
	).Scan(&h.ID, &h.Name, &h.TotalSeats, &h.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return h, err
}

func (r *HallRepo) GetSeatMap(ctx context.Context, hallID uuid.UUID) ([]*domain.Seat, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,hall_id,row_number,seat_number,seat_type FROM seats WHERE hall_id=$1 ORDER BY row_number,seat_number`,
		hallID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var seats []*domain.Seat
	for rows.Next() {
		s := &domain.Seat{}
		if err := rows.Scan(&s.ID, &s.HallID, &s.RowNumber, &s.SeatNumber, &s.SeatType); err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}
	return seats, nil
}

// ── Showtime Repository ──────────────────────────────────
type ShowtimeRepo struct{ db *pgxpool.Pool }

func NewShowtimeRepo(db *pgxpool.Pool) *ShowtimeRepo { return &ShowtimeRepo{db: db} }

func (r *ShowtimeRepo) Create(ctx context.Context, s *domain.Showtime) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO showtimes (id,movie_id,hall_id,start_time,end_time,price) VALUES ($1,$2,$3,$4,$5,$6)`,
		s.ID, s.MovieID, s.HallID, s.StartTime, s.EndTime, s.Price,
	)
	return err
}

func (r *ShowtimeRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Showtime, error) {
	s := &domain.Showtime{}
	err := r.db.QueryRow(ctx,
		`SELECT id,movie_id,hall_id,start_time,end_time,price FROM showtimes WHERE id=$1`, id,
	).Scan(&s.ID, &s.MovieID, &s.HallID, &s.StartTime, &s.EndTime, &s.Price)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return s, err
}

func (r *ShowtimeRepo) List(ctx context.Context, movieID uuid.UUID, date string) ([]*domain.Showtime, error) {
	query := `SELECT id,movie_id,hall_id,start_time,end_time,price FROM showtimes WHERE movie_id=$1`
	args := []interface{}{movieID}
	if date != "" {
		query += ` AND DATE(start_time)=$2`
		args = append(args, date)
	}
	query += ` ORDER BY start_time`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var showtimes []*domain.Showtime
	for rows.Next() {
		s := &domain.Showtime{}
		if err := rows.Scan(&s.ID, &s.MovieID, &s.HallID, &s.StartTime, &s.EndTime, &s.Price); err != nil {
			return nil, err
		}
		showtimes = append(showtimes, s)
	}
	return showtimes, nil
}

func (r *ShowtimeRepo) GetAvailableSeats(ctx context.Context, showtimeID uuid.UUID) ([]*domain.Seat, error) {
	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.hall_id, s.row_number, s.seat_number, s.seat_type
		FROM seats s
		JOIN showtimes st ON st.hall_id = s.hall_id
		WHERE st.id = $1
		  AND s.id NOT IN (
		    SELECT bs.seat_id FROM booking_seats bs
		    JOIN bookings b ON b.id = bs.booking_id
		    WHERE b.showtime_id = $1 AND b.status = 'confirmed'
		  )
		ORDER BY s.row_number, s.seat_number`,
		showtimeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var seats []*domain.Seat
	for rows.Next() {
		s := &domain.Seat{}
		if err := rows.Scan(&s.ID, &s.HallID, &s.RowNumber, &s.SeatNumber, &s.SeatType); err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}
	return seats, nil
}

// ── Booking Repository ───────────────────────────────────
type BookingRepo struct{ db *pgxpool.Pool }

func NewBookingRepo(db *pgxpool.Pool) *BookingRepo { return &BookingRepo{db: db} }

func (r *BookingRepo) CreateWithSeats(ctx context.Context, b *domain.Booking) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Lock seat rows to prevent concurrent booking
	_, err = tx.Exec(ctx,
		`SELECT id FROM seats WHERE id = ANY($1) FOR UPDATE`, b.SeatIDs,
	)
	if err != nil {
		return err
	}

	// Check none are already booked for this showtime
	var taken int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*) FROM booking_seats bs
		JOIN bookings bk ON bk.id = bs.booking_id
		WHERE bs.seat_id = ANY($1)
		  AND bk.showtime_id = $2
		  AND bk.status = 'confirmed'`,
		b.SeatIDs, b.ShowtimeID,
	).Scan(&taken)
	if err != nil {
		return err
	}
	if taken > 0 {
		return domain.ErrSeatsUnavailable
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO bookings (id,user_id,showtime_id,status,total_price,booked_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		b.ID, b.UserID, b.ShowtimeID, b.Status, b.TotalPrice, b.BookedAt,
	)
	if err != nil {
		return err
	}

	for _, seatID := range b.SeatIDs {
		_, err = tx.Exec(ctx,
			`INSERT INTO booking_seats (id,booking_id,seat_id) VALUES ($1,$2,$3)`,
			uuid.New(), b.ID, seatID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *BookingRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	b := &domain.Booking{}
	err := r.db.QueryRow(ctx,
		`SELECT id,user_id,showtime_id,status,total_price,booked_at FROM bookings WHERE id=$1`, id,
	).Scan(&b.ID, &b.UserID, &b.ShowtimeID, &b.Status, &b.TotalPrice, &b.BookedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	b.SeatIDs, err = r.getSeatIDs(ctx, id)
	return b, err
}

func (r *BookingRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Booking, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id,user_id,showtime_id,status,total_price,booked_at FROM bookings WHERE user_id=$1 ORDER BY booked_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var bookings []*domain.Booking
	for rows.Next() {
		b := &domain.Booking{}
		if err := rows.Scan(&b.ID, &b.UserID, &b.ShowtimeID, &b.Status, &b.TotalPrice, &b.BookedAt); err != nil {
			return nil, err
		}
		b.SeatIDs, _ = r.getSeatIDs(ctx, b.ID)
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *BookingRepo) Cancel(ctx context.Context, bookingID, userID uuid.UUID) (*domain.Booking, error) {
	now := time.Now()
	var b domain.Booking
	err := r.db.QueryRow(ctx,
		`UPDATE bookings SET status='cancelled', cancelled_at=$1
		 WHERE id=$2 AND user_id=$3 AND status='confirmed'
		 RETURNING id,user_id,showtime_id,status,total_price,booked_at`,
		now, bookingID, userID,
	).Scan(&b.ID, &b.UserID, &b.ShowtimeID, &b.Status, &b.TotalPrice, &b.BookedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("%w: booking not found or already cancelled", domain.ErrNotFound)
	}
	return &b, err
}

func (r *BookingRepo) getSeatIDs(ctx context.Context, bookingID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `SELECT seat_id FROM booking_seats WHERE booking_id=$1`, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
