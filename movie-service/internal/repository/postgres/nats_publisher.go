package postgres

import (
	"context"
	"encoding/json"

	"cinema/movie-service/internal/domain"
	"github.com/nats-io/nats.go"
)

const (
	SubjectBookingConfirmed = "booking.confirmed"
	SubjectBookingCancelled = "booking.cancelled"
)

type NATSPublisher struct {
	nc *nats.Conn
}

func NewNATSPublisher(nc *nats.Conn) *NATSPublisher {
	return &NATSPublisher{nc: nc}
}

func (p *NATSPublisher) PublishBookingConfirmed(_ context.Context, event *domain.BookingConfirmedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.nc.Publish(SubjectBookingConfirmed, data)
}

func (p *NATSPublisher) PublishBookingCancelled(_ context.Context, event *domain.BookingCancelledEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.nc.Publish(SubjectBookingCancelled, data)
}
