FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates

# Copy everything first so go mod tidy can see all imports
COPY . .

# Tidy resolves all deps (including proto subpackages) and writes go.sum
RUN go mod tidy

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/payment-service .
COPY --from=builder /app/migrations ./migrations
EXPOSE 50053
CMD ["./payment-service"]
