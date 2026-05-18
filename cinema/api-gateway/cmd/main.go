package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cinema/api-gateway/internal/router"
	"cinema/api-gateway/internal/telemetry"
)

func main() {
	movieAddr := getEnv("MOVIE_SERVICE_ADDR", "localhost:50051")
	authAddr := getEnv("AUTH_SERVICE_ADDR", "localhost:50052")
	paymentAddr := getEnv("PAYMENT_SERVICE_ADDR", "localhost:50053")
	port := getEnv("HTTP_PORT", "8080")

	ctx := context.Background()
	shutdown, err := telemetry.Setup(ctx, getEnv("OTEL_SERVICE_NAME", "api-gateway"), getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "tempo:4318"))
	if err != nil {
		log.Printf("warning: telemetry setup failed: %v", err)
	} else {
		defer func() {
			if err := shutdown(ctx); err != nil {
				log.Printf("failed to shutdown telemetry: %v", err)
			}
		}()
	}

	gw, err := router.NewGateway(movieAddr, authAddr, paymentAddr)
	if err != nil {
		log.Fatalf("failed to connect to services: %v", err)
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("api-gateway listening on %s", addr)
	log.Printf("  -> movie-service  : %s", movieAddr)
	log.Printf("  -> auth-service   : %s", authAddr)
	log.Printf("  -> payment-service: %s", paymentAddr)

	handler := telemetry.WrapHTTPHandler(gw.Router(), "api-gateway")
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
