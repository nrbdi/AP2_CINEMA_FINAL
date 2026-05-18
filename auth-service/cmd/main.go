package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	grpcDelivery "cinema/auth-service/internal/delivery/grpc"
	emailPkg "cinema/auth-service/internal/pkg/email"
	repoPg "cinema/auth-service/internal/repository/postgres"
	"cinema/auth-service/internal/usecase"
	pb "cinema/proto/auth"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	dbURL := getEnv("DB_URL", "postgres://auth_user:auth_pass@localhost:5433/cinema_auth?sslmode=disable")
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("postgres connect failed: %v", err)
	}
	defer pool.Close()
	log.Println("connected to postgres")

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("migration init: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration: %v", err)
	}
	log.Println("migrations applied")

	redisURL := getEnv("REDIS_URL", "localhost:6380")
	rdb := redis.NewClient(&redis.Options{Addr: redisURL})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("redis warning: %v", err)
	}
	log.Println("connected to redis")

	userRepo    := repoPg.NewUserRepo(pool)
	tokenCache  := repoPg.NewTokenCache(rdb)
	emailSender := emailPkg.NewSMTPSender()

	uc      := usecase.NewAuthUsecase(userRepo, tokenCache, emailSender)
	handler := grpcDelivery.NewAuthHandler(uc)

	port := getEnv("GRPC_PORT", "50052")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterAuthServiceServer(srv, handler)
	reflection.Register(srv)

	log.Printf("auth-service listening on :%s", port)
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
