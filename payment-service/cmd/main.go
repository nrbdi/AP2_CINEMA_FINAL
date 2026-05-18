package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"cinema/payment-service/internal/consumer"
	grpcDelivery "cinema/payment-service/internal/delivery/grpc"
	emailPkg "cinema/payment-service/internal/pkg/email"
	repoPg "cinema/payment-service/internal/repository/postgres"
	"cinema/payment-service/internal/usecase"
	pb "cinema/proto/payment"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	dbURL := getEnv("DB_URL", "postgres://payment_user:payment_pass@localhost:5434/cinema_payment?sslmode=disable")
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

	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer nc.Close()
	log.Println("connected to nats")

	paymentRepo := repoPg.NewPaymentRepo(pool)
	paymentUC   := usecase.NewPaymentUsecase(paymentRepo)
	emailSender := emailPkg.NewSMTPSender()

	natsCons := consumer.NewNATSConsumer(nc, paymentUC, emailSender)
	if err := natsCons.Start(); err != nil {
		log.Fatalf("nats consumer start: %v", err)
	}

	handler := grpcDelivery.NewPaymentHandler(paymentUC)
	port := getEnv("GRPC_PORT", "50053")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterPaymentServiceServer(srv, handler)
	reflection.Register(srv)

	log.Printf("payment-service listening on :%s", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
