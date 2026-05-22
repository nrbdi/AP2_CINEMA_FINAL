# Cinema — Compact README

A compact overview of the Cinema microservices project and how to run it. This README keeps only the essential information for developers and operators.

## What this repo contains
- API Gateway (HTTP) that routes to gRPC services
- Movie service (movies, halls, showtimes, bookings)
- Auth service (users, auth tokens)
- Payment service (payments, receipts, email notifications)
- Docker Compose for local development and observability (Prometheus, Grafana, Loki, Tempo)

## Quick Start (local)
1. Ensure Docker and Docker Compose are installed.
2. From the `cinema` folder run:

```bash
docker compose up -d
```

3. To view logs:

```bash
docker compose logs -f api-gateway
```

4. Stop and remove containers:

```bash
docker compose down
```

## Configuration notes
- Environment variables are read from the Compose file and `./.env` (if present).
- For SMTP (email) set `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, and `SMTP_FROM` in `./.env` and restart the services.
- Grafana is available on port `3000` by default when the monitoring stack is up.

## Main HTTP endpoints (examples)
- POST `/api/v1/auth/register` — create account (sends welcome email)
- POST `/api/v1/auth/login` — obtain tokens
- GET `/api/v1/movies` — list movies
- GET `/api/v1/movies/{id}` — movie details + showtimes
- POST `/api/v1/bookings` — create a booking (requires auth)
- GET `/api/v1/payments/{id}/receipt` — download receipt

For a full API reference, see the `proto/` definitions and the gateway router.

## Endpoint Assignment (team ownership)
- **Aknur — Auth & Payment (12 endpoints)**: registration, login, token flows, profile management, payment create/list/receipt.
- **Dinara — Movie (12 endpoints)**: movie CRUD, halls, showtimes, seat availability, booking create/get.

Notes: `ListUserBookings` and `CancelBooking` are integration points; assign during integration testing if needed.

## Tests
- Unit tests live under each service `internal/usecase` folder. Run service tests locally with `go test ./...` inside the service folder.

## Project layout (short)

```
cinema/
├── docker-compose.yml
├── proto/
├── api-gateway/
├── movie-service/
├── auth-service/
└── payment-service/
```

---



