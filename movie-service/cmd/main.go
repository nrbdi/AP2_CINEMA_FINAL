package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	grpcDelivery "cinema/movie-service/internal/delivery/grpc"
	repoPostgres "cinema/movie-service/internal/repository/postgres"
	"cinema/movie-service/internal/usecase"
	pb "cinema/movie-service/proto/movie"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	dbURL := getEnv("DB_URL", "postgres://movie_user:movie_pass@localhost:5432/cinema_movie?sslmode=disable")
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("postgres ping failed: %v", err)
	}
	log.Println("connected to postgres")

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("migration init failed: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations applied")

	redisURL := getEnv("REDIS_URL", "localhost:6379")
	rdb := redis.NewClient(&redis.Options{Addr: redisURL})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("redis unavailable: %v", err)
	}
	log.Println("connected to redis")

	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("nats connection failed: %v", err)
	}
	defer nc.Close()
	log.Println("connected to nats")

	movieRepo := repoPostgres.NewMovieRepo(pool)
	hallRepo := repoPostgres.NewHallRepo(pool)
	showtimeRepo := repoPostgres.NewShowtimeRepo(pool)
	bookingRepo := repoPostgres.NewBookingRepo(pool)
	publisher := repoPostgres.NewNATSPublisher(nc)
	seatCache := repoPostgres.NewSeatCache(rdb)

	uc := usecase.NewMovieUsecase(movieRepo, hallRepo, showtimeRepo, bookingRepo, publisher, seatCache)
	handler := grpcDelivery.NewMovieHandler(uc)

	port := getEnv("GRPC_PORT", "50051")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterMovieServiceServer(srv, handler)
	reflection.Register(srv)

	log.Printf("movie-service listening on :%s", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
