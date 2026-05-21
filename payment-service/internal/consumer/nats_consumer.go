package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cinema/payment-service/internal/domain"
	"cinema/payment-service/internal/usecase"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

const SubjectBookingConfirmed = "booking.confirmed"

type EmailSender interface {
	SendPaymentReceipt(toEmail, toName, movieTitle, showtime string, seats []string, amount float64) error
}

type NATSConsumer struct {
	nc          *nats.Conn
	paymentUC   *usecase.PaymentUsecase
	emailSender EmailSender
}

func NewNATSConsumer(nc *nats.Conn, uc *usecase.PaymentUsecase, email EmailSender) *NATSConsumer {
	return &NATSConsumer{nc: nc, paymentUC: uc, emailSender: email}
}

func (c *NATSConsumer) Start() error {
	_, err := c.nc.Subscribe(SubjectBookingConfirmed, func(msg *nats.Msg) {
		var event domain.BookingConfirmedEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("failed to unmarshal booking event: %v", err)
			return
		}
		c.handleBookingConfirmed(event)
	})
	if err != nil {
		return fmt.Errorf("nats subscribe error: %w", err)
	}
	log.Printf("subscribed to %s", SubjectBookingConfirmed)
	return nil
}

func (c *NATSConsumer) handleBookingConfirmed(event domain.BookingConfirmedEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bookingID, err := uuid.Parse(event.BookingID)
	if err != nil {
		log.Printf("invalid booking_id: %v", err)
		return
	}
	userID, err := uuid.Parse(event.UserID)
	if err != nil {
		log.Printf("invalid user_id: %v", err)
		return
	}

	payment, err := c.paymentUC.CreatePayment(ctx, bookingID, userID, event.TotalPrice, "system")
	if err != nil {
		log.Printf("failed to create payment for booking %s: %v", event.BookingID, err)
		return
	}

	log.Printf("payment %s created for booking %s amount=%.2f", payment.ID, event.BookingID, payment.Amount)

	
	seatLabels := make([]string, len(event.SeatIDs))
	for i, id := range event.SeatIDs {
		seatLabels[i] = id[:8] 
	}
	go func() {
		if err := c.emailSender.SendPaymentReceipt(
			"user@example.com", 
			"Customer",
			"Movie Title",
			event.ShowtimeID,
			seatLabels,
			event.TotalPrice,
		); err != nil {
			log.Printf("email send failed: %v", err)
		}
	}()
}
