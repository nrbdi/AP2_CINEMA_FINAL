package router

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"cinema/api-gateway/internal/middleware"
	"cinema/api-gateway/internal/telemetry"
	pbAuth "cinema/proto/auth"
	pbMovie "cinema/proto/movie"
	pbPayment "cinema/proto/payment"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	movieClient   pbMovie.MovieServiceClient
	authClient    pbAuth.AuthServiceClient
	paymentClient pbPayment.PaymentServiceClient
}

func NewGateway(movieAddr, authAddr, paymentAddr string) (*Gateway, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	dialOpts = append(dialOpts, telemetry.OTelGRPCDialOptions()...)

	movieConn, err := grpc.Dial(movieAddr, dialOpts...)
	if err != nil {
		return nil, err
	}
	authConn, err := grpc.Dial(authAddr, dialOpts...)
	if err != nil {
		return nil, err
	}
	paymentConn, err := grpc.Dial(paymentAddr, dialOpts...)
	if err != nil {
		return nil, err
	}
	return &Gateway{
		movieClient:   pbMovie.NewMovieServiceClient(movieConn),
		authClient:    pbAuth.NewAuthServiceClient(authConn),
		paymentClient: pbPayment.NewPaymentServiceClient(paymentConn),
	}, nil
}

func (g *Gateway) Router() http.Handler {
	r := mux.NewRouter()
	r.Use(middleware.CORS)
	r.Use(middleware.Metrics)
	r.Use(middleware.Logger)
	r.Use(middleware.RateLimit(100))

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		jsonResp(w, map[string]string{"status": "ok"})
	}).Methods("GET")

	r.Handle("/metrics", promhttp.Handler())

	r.HandleFunc("/api/v1/auth/register", g.register).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", g.login).Methods("POST")
	r.HandleFunc("/api/v1/auth/refresh", g.refreshToken).Methods("POST")

	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.Auth)
	protected.HandleFunc("/auth/logout", g.logout).Methods("POST")
	protected.HandleFunc("/auth/profile", g.getProfile).Methods("GET")
	protected.HandleFunc("/auth/profile", g.updateProfile).Methods("PUT")

	r.HandleFunc("/api/v1/movies", g.listMovies).Methods("GET")
	r.HandleFunc("/api/v1/movies/{id}", g.getMovie).Methods("GET")
	r.HandleFunc("/api/v1/movies/{id}/showtimes", g.listShowtimes).Methods("GET")
	r.HandleFunc("/api/v1/showtimes/{id}/seats", g.getAvailableSeats).Methods("GET")

	admin := r.PathPrefix("/api/v1/admin").Subrouter()
	admin.Use(middleware.Auth)
	admin.Use(middleware.AdminOnly)
	admin.HandleFunc("/movies", g.createMovie).Methods("POST")
	admin.HandleFunc("/movies/{id}", g.updateMovie).Methods("PUT")
	admin.HandleFunc("/movies/{id}", g.deleteMovie).Methods("DELETE")
	admin.HandleFunc("/halls", g.createHall).Methods("POST")
	admin.HandleFunc("/halls/{id}/seats", g.getHallSeatMap).Methods("GET")
	admin.HandleFunc("/showtimes", g.createShowtime).Methods("POST")

	protected.HandleFunc("/bookings", g.createBooking).Methods("POST")
	protected.HandleFunc("/bookings", g.listUserBookings).Methods("GET")
	protected.HandleFunc("/bookings/{id}", g.getBooking).Methods("GET")
	protected.HandleFunc("/bookings/{id}/cancel", g.cancelBooking).Methods("POST")

	protected.HandleFunc("/payments", g.listUserPayments).Methods("GET")
	protected.HandleFunc("/payments/{id}", g.getPayment).Methods("GET")
	protected.HandleFunc("/payments/{id}/receipt", g.getReceipt).Methods("GET")

	return r
}

func rpcCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func jsonResp(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func jsonErr(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (g *Gateway) register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		Phone    string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonErr(w, 400, "invalid body")
		return
	}
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.authClient.Register(ctx, &pbAuth.RegisterRequest{
		Email: req.Email, Password: req.Password, FullName: req.FullName, Phone: req.Phone,
	})
	if err != nil {
		jsonErr(w, 400, err.Error())
		return
	}
	w.WriteHeader(201)
	jsonResp(w, resp)
}

func (g *Gateway) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.authClient.Login(ctx, &pbAuth.LoginRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		jsonErr(w, 401, "invalid credentials")
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if len(token) > 7 {
		token = token[7:]
	}
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.authClient.Logout(ctx, &pbAuth.LogoutRequest{AccessToken: token})
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) refreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.authClient.RefreshToken(ctx, &pbAuth.RefreshTokenRequest{RefreshToken: req.RefreshToken})
	if err != nil {
		jsonErr(w, 401, "invalid refresh token")
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) getProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.authClient.GetProfile(ctx, &pbAuth.GetProfileRequest{UserId: userID})
	if err != nil {
		jsonErr(w, 404, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) updateProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	var req struct {
		FullName string `json:"full_name"`
		Phone    string `json:"phone"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.authClient.UpdateProfile(ctx, &pbAuth.UpdateProfileRequest{
		UserId: userID, FullName: req.FullName, Phone: req.Phone,
	})
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) listMovies(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.ListMovies(ctx, &pbMovie.ListMoviesRequest{Genre: genre, Page: 1, PageSize: 20})
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) getMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.GetMovie(ctx, &pbMovie.GetMovieRequest{Id: id})
	if err != nil {
		jsonErr(w, 404, "not found")
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) createMovie(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DurationMin int32  `json:"duration_min"`
		Genre       string `json:"genre"`
		PosterUrl   string `json:"poster_url"`
		ReleaseDate string `json:"release_date"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.CreateMovie(ctx, &pbMovie.CreateMovieRequest{
		Title: req.Title, Description: req.Description, DurationMin: req.DurationMin,
		Genre: req.Genre, PosterUrl: req.PosterUrl, ReleaseDate: req.ReleaseDate,
	})
	if err != nil {
		jsonErr(w, 400, err.Error())
		return
	}
	w.WriteHeader(201)
	jsonResp(w, resp)
}

func (g *Gateway) updateMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DurationMin int32  `json:"duration_min"`
		Genre       string `json:"genre"`
		PosterUrl   string `json:"poster_url"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.UpdateMovie(ctx, &pbMovie.UpdateMovieRequest{
		Id: id, Title: req.Title, Description: req.Description,
		DurationMin: req.DurationMin, Genre: req.Genre, PosterUrl: req.PosterUrl,
	})
	if err != nil {
		jsonErr(w, 400, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) deleteMovie(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.DeleteMovie(ctx, &pbMovie.DeleteMovieRequest{Id: id})
	if err != nil {
		jsonErr(w, 400, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) createHall(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Rows        int32  `json:"rows"`
		SeatsPerRow int32  `json:"seats_per_row"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.CreateHall(ctx, &pbMovie.CreateHallRequest{
		Name: req.Name, Rows: req.Rows, SeatsPerRow: req.SeatsPerRow,
	})
	if err != nil {
		jsonErr(w, 400, err.Error())
		return
	}
	w.WriteHeader(201)
	jsonResp(w, resp)
}

func (g *Gateway) getHallSeatMap(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.GetHallSeatMap(ctx, &pbMovie.GetHallRequest{HallId: id})
	if err != nil {
		jsonErr(w, 404, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) createShowtime(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MovieId   string  `json:"movie_id"`
		HallId    string  `json:"hall_id"`
		StartTime string  `json:"start_time"`
		EndTime   string  `json:"end_time"`
		Price     float64 `json:"price"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.CreateShowtime(ctx, &pbMovie.CreateShowtimeRequest{
		MovieId: req.MovieId, HallId: req.HallId,
		StartTime: req.StartTime, EndTime: req.EndTime, Price: req.Price,
	})
	if err != nil {
		jsonErr(w, 400, err.Error())
		return
	}
	w.WriteHeader(201)
	jsonResp(w, resp)
}

func (g *Gateway) listShowtimes(w http.ResponseWriter, r *http.Request) {
	movieID := mux.Vars(r)["id"]
	date := r.URL.Query().Get("date")
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.ListShowtimes(ctx, &pbMovie.ListShowtimesRequest{MovieId: movieID, Date: date})
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) getAvailableSeats(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.GetAvailableSeats(ctx, &pbMovie.GetAvailableSeatsRequest{ShowtimeId: id})
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) createBooking(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	var req struct {
		ShowtimeId string   `json:"showtime_id"`
		SeatIds    []string `json:"seat_ids"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.CreateBooking(ctx, &pbMovie.CreateBookingRequest{
		UserId: userID, ShowtimeId: req.ShowtimeId, SeatIds: req.SeatIds,
	})
	if err != nil {
		jsonErr(w, 400, err.Error())
		return
	}
	w.WriteHeader(201)
	jsonResp(w, resp)
}

func (g *Gateway) getBooking(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.GetBooking(ctx, &pbMovie.GetBookingRequest{BookingId: id})
	if err != nil {
		jsonErr(w, 404, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) listUserBookings(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.ListUserBookings(ctx, &pbMovie.ListUserBookingsRequest{UserId: userID})
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) cancelBooking(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	id := mux.Vars(r)["id"]
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.movieClient.CancelBooking(ctx, &pbMovie.CancelBookingRequest{BookingId: id, UserId: userID})
	if err != nil {
		jsonErr(w, 400, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) listUserPayments(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.ContextUserID).(string)
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.paymentClient.ListUserPayments(ctx, &pbPayment.ListUserPaymentsRequest{UserId: userID})
	if err != nil {
		jsonErr(w, 500, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) getPayment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.paymentClient.GetPayment(ctx, &pbPayment.GetPaymentRequest{PaymentId: id})
	if err != nil {
		jsonErr(w, 404, err.Error())
		return
	}
	jsonResp(w, resp)
}

func (g *Gateway) getReceipt(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx, cancel := rpcCtx()
	defer cancel()
	resp, err := g.paymentClient.GetReceipt(ctx, &pbPayment.GetReceiptRequest{PaymentId: id})
	if err != nil {
		jsonErr(w, 404, err.Error())
		return
	}
	jsonResp(w, resp)
}
