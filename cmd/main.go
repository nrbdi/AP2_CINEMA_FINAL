package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"cinema/api-gateway/internal/router"
)

func main() {
	movieAddr := getEnv("MOVIE_SERVICE_ADDR", "localhost:50051")
	authAddr := getEnv("AUTH_SERVICE_ADDR", "localhost:50052")
	paymentAddr := getEnv("PAYMENT_SERVICE_ADDR", "localhost:50053")
	port := getEnv("HTTP_PORT", "8080")

	gw, err := router.NewGateway(movieAddr, authAddr, paymentAddr)
	if err != nil {
		log.Fatalf("failed to connect to services: %v", err)
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("api-gateway listening on %s", addr)
	log.Printf("  -> movie-service  : %s", movieAddr)
	log.Printf("  -> auth-service   : %s", authAddr)
	log.Printf("  -> payment-service: %s", paymentAddr)

	if err := http.ListenAndServe(addr, gw.Router()); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
