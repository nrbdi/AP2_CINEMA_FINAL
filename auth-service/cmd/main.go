package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	emailPkg "cinema/auth-service/internal/pkg/email"
	"cinema/auth-service/internal/delivery/grpc"
	repoPg "cinema/auth-service/internal/repository"
	"cinema/auth-service/internal/usecase"
	pb "cinema/auth-service/proto/auth"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	googleGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://auth_user:auth_pass@localhost:5433/cinema_auth?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}
	defer pool.Close()

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("migration init failed: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations applied")

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6380"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisURL})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("redis warning: %v", err)
	}

	userRepo   := repoPg.NewUserRepo(pool)
	tokenCache := repoPg.NewTokenCache(rdb)
	emailSender := emailPkg.NewSMTPSender()

	uc      := usecase.NewAuthUsecase(userRepo, tokenCache, emailSender)
	handler := grpc.NewAuthHandler(uc)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	srv := googleGRPC.NewServer()
	pb.RegisterAuthServiceServer(srv, handler)
	reflection.Register(srv)

	log.Printf("auth-service listening on :%s", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve failed: %v", err)
	}
}
