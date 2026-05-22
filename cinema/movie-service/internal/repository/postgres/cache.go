package postgres

import (
	"context"
	"encoding/json"
	"time"

	"cinema/movie-service/internal/domain"
	"github.com/redis/go-redis/v9"
)

type SeatCacheRedis struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSeatCache(client *redis.Client) *SeatCacheRedis {
	return &SeatCacheRedis{client: client, ttl: 5 * time.Minute}
}

func cacheKey(showtimeID string) string {
	return "seats:available:" + showtimeID
}

func (c *SeatCacheRedis) SetAvailableSeats(ctx context.Context, showtimeID string, seats []*domain.Seat) error {
	data, err := json.Marshal(seats)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, cacheKey(showtimeID), data, c.ttl).Err()
}

func (c *SeatCacheRedis) GetAvailableSeats(ctx context.Context, showtimeID string) ([]*domain.Seat, bool, error) {
	data, err := c.client.Get(ctx, cacheKey(showtimeID)).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var seats []*domain.Seat
	if err := json.Unmarshal(data, &seats); err != nil {
		return nil, false, err
	}
	return seats, true, nil
}

func (c *SeatCacheRedis) InvalidateShowtime(ctx context.Context, showtimeID string) error {
	return c.client.Del(ctx, cacheKey(showtimeID)).Err()
}
