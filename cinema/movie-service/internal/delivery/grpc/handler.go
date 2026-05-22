package grpc

import (
	"context"

	"cinema/movie-service/internal/domain"
	"cinema/movie-service/internal/usecase"
	pb "cinema/proto/movie"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MovieHandler struct {
	pb.UnimplementedMovieServiceServer
	uc *usecase.MovieUsecase
}

func NewMovieHandler(uc *usecase.MovieUsecase) *MovieHandler {
	return &MovieHandler{uc: uc}
}

func mapErr(err error) error {
	switch err {
	case domain.ErrNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrSeatsUnavailable:
		return status.Error(codes.AlreadyExists, err.Error())
	case domain.ErrUnauthorized:
		return status.Error(codes.PermissionDenied, err.Error())
	case domain.ErrInvalidInput:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func movieToPB(m *domain.Movie) *pb.Movie {
	return &pb.Movie{
		Id:          m.ID.String(),
		Title:       m.Title,
		Description: m.Description,
		DurationMin: int32(m.DurationMin),
		Genre:       m.Genre,
		PosterUrl:   m.PosterURL,
		ReleaseDate: m.ReleaseDate,
		CreatedAt:   timestamppb.New(m.CreatedAt),
	}
}

func seatToPB(s *domain.Seat) *pb.Seat {
	return &pb.Seat{
		Id:         s.ID.String(),
		HallId:     s.HallID.String(),
		RowNumber:  int32(s.RowNumber),
		SeatNumber: int32(s.SeatNumber),
		SeatType:   string(s.SeatType),
	}
}

func bookingToPB(b *domain.Booking) *pb.Booking {
	seatIDs := make([]string, len(b.SeatIDs))
	for i, id := range b.SeatIDs {
		seatIDs[i] = id.String()
	}
	return &pb.Booking{
		Id:         b.ID.String(),
		UserId:     b.UserID.String(),
		ShowtimeId: b.ShowtimeID.String(),
		Status:     string(b.Status),
		SeatIds:    seatIDs,
		TotalPrice: b.TotalPrice,
		BookedAt:   timestamppb.New(b.BookedAt),
	}
}

// ── Movie endpoints ──────────────────────────────────────
func (h *MovieHandler) CreateMovie(ctx context.Context, req *pb.CreateMovieRequest) (*pb.MovieResponse, error) {
	m, err := h.uc.CreateMovie(ctx, &domain.Movie{
		Title:       req.Title,
		Description: req.Description,
		DurationMin: int(req.DurationMin),
		Genre:       req.Genre,
		PosterURL:   req.PosterUrl,
		ReleaseDate: req.ReleaseDate,
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.MovieResponse{Movie: movieToPB(m)}, nil
}

func (h *MovieHandler) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.MovieResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	m, err := h.uc.GetMovie(ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.MovieResponse{Movie: movieToPB(m)}, nil
}

func (h *MovieHandler) ListMovies(ctx context.Context, req *pb.ListMoviesRequest) (*pb.ListMoviesResponse, error) {
	movies, total, err := h.uc.ListMovies(ctx, req.Genre, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, mapErr(err)
	}
	pbMovies := make([]*pb.Movie, len(movies))
	for i, m := range movies {
		pbMovies[i] = movieToPB(m)
	}
	return &pb.ListMoviesResponse{Movies: pbMovies, Total: int32(total)}, nil
}

func (h *MovieHandler) UpdateMovie(ctx context.Context, req *pb.UpdateMovieRequest) (*pb.MovieResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	m, err := h.uc.UpdateMovie(ctx, &domain.Movie{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		DurationMin: int(req.DurationMin),
		Genre:       req.Genre,
		PosterURL:   req.PosterUrl,
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.MovieResponse{Movie: movieToPB(m)}, nil
}

func (h *MovieHandler) DeleteMovie(ctx context.Context, req *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	if err := h.uc.DeleteMovie(ctx, id); err != nil {
		return nil, mapErr(err)
	}
	return &pb.DeleteMovieResponse{Success: true}, nil
}

// ── Hall endpoints ───────────────────────────────────────
func (h *MovieHandler) CreateHall(ctx context.Context, req *pb.CreateHallRequest) (*pb.HallResponse, error) {
	hall, err := h.uc.CreateHall(ctx, req.Name, int(req.Rows), int(req.SeatsPerRow))
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.HallResponse{Hall: &pb.Hall{
		Id: hall.ID.String(), Name: hall.Name, TotalSeats: int32(hall.TotalSeats),
	}}, nil
}

func (h *MovieHandler) GetHallSeatMap(ctx context.Context, req *pb.GetHallRequest) (*pb.SeatMapResponse, error) {
	hallID, err := uuid.Parse(req.HallId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid hall_id")
	}
	seats, err := h.uc.GetHallSeatMap(ctx, hallID)
	if err != nil {
		return nil, mapErr(err)
	}
	pbSeats := make([]*pb.Seat, len(seats))
	for i, s := range seats {
		pbSeats[i] = seatToPB(s)
	}
	return &pb.SeatMapResponse{Seats: pbSeats}, nil
}

// ── Showtime endpoints ───────────────────────────────────
func (h *MovieHandler) CreateShowtime(ctx context.Context, req *pb.CreateShowtimeRequest) (*pb.ShowtimeResponse, error) {
	movieID, _ := uuid.Parse(req.MovieId)
	hallID, _ := uuid.Parse(req.HallId)

	st, err := h.uc.CreateShowtime(ctx, &domain.Showtime{
		MovieID: movieID,
		HallID:  hallID,
		Price:   req.Price,
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.ShowtimeResponse{Showtime: &pb.Showtime{
		Id:        st.ID.String(),
		MovieId:   st.MovieID.String(),
		HallId:    st.HallID.String(),
		StartTime: st.StartTime.Format("2006-01-02T15:04:05Z"),
		EndTime:   st.EndTime.Format("2006-01-02T15:04:05Z"),
		Price:     st.Price,
	}}, nil
}

func (h *MovieHandler) ListShowtimes(ctx context.Context, req *pb.ListShowtimesRequest) (*pb.ListShowtimesResponse, error) {
	movieID, _ := uuid.Parse(req.MovieId)
	showtimes, err := h.uc.ListShowtimes(ctx, movieID, req.Date)
	if err != nil {
		return nil, mapErr(err)
	}
	pbST := make([]*pb.Showtime, len(showtimes))
	for i, s := range showtimes {
		pbST[i] = &pb.Showtime{
			Id: s.ID.String(), MovieId: s.MovieID.String(), HallId: s.HallID.String(),
			StartTime: s.StartTime.Format("2006-01-02T15:04:05Z"),
			EndTime:   s.EndTime.Format("2006-01-02T15:04:05Z"),
			Price:     s.Price,
		}
	}
	return &pb.ListShowtimesResponse{Showtimes: pbST}, nil
}

func (h *MovieHandler) GetAvailableSeats(ctx context.Context, req *pb.GetAvailableSeatsRequest) (*pb.AvailableSeatsResponse, error) {
	id, _ := uuid.Parse(req.ShowtimeId)
	seats, err := h.uc.GetAvailableSeats(ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	pbSeats := make([]*pb.Seat, len(seats))
	for i, s := range seats {
		pbSeats[i] = seatToPB(s)
	}
	return &pb.AvailableSeatsResponse{Seats: pbSeats}, nil
}

// ── Booking endpoints ────────────────────────────────────
func (h *MovieHandler) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.BookingResponse, error) {
	userID, _ := uuid.Parse(req.UserId)
	showtimeID, _ := uuid.Parse(req.ShowtimeId)
	seatIDs := make([]uuid.UUID, len(req.SeatIds))
	for i, s := range req.SeatIds {
		seatIDs[i], _ = uuid.Parse(s)
	}
	b, err := h.uc.CreateBooking(ctx, userID, showtimeID, seatIDs)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.BookingResponse{Booking: bookingToPB(b)}, nil
}

func (h *MovieHandler) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.BookingResponse, error) {
	id, _ := uuid.Parse(req.BookingId)
	b, err := h.uc.GetBooking(ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.BookingResponse{Booking: bookingToPB(b)}, nil
}

func (h *MovieHandler) ListUserBookings(ctx context.Context, req *pb.ListUserBookingsRequest) (*pb.ListBookingsResponse, error) {
	userID, _ := uuid.Parse(req.UserId)
	bookings, err := h.uc.ListUserBookings(ctx, userID)
	if err != nil {
		return nil, mapErr(err)
	}
	pbBookings := make([]*pb.Booking, len(bookings))
	for i, b := range bookings {
		pbBookings[i] = bookingToPB(b)
	}
	return &pb.ListBookingsResponse{Bookings: pbBookings}, nil
}

func (h *MovieHandler) CancelBooking(ctx context.Context, req *pb.CancelBookingRequest) (*pb.BookingResponse, error) {
	bookingID, _ := uuid.Parse(req.BookingId)
	userID, _ := uuid.Parse(req.UserId)
	b, err := h.uc.CancelBooking(ctx, bookingID, userID)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.BookingResponse{Booking: bookingToPB(b)}, nil
}
