package usecase

import (
	"context"
	"time"

	"cinema/payment-service/internal/domain"
	"github.com/google/uuid"
)

type PaymentUsecase struct {
	repo domain.PaymentRepository
}

func NewPaymentUsecase(repo domain.PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{repo: repo}
}

func (uc *PaymentUsecase) CreatePayment(ctx context.Context, bookingID, userID uuid.UUID, amount float64, method string) (*domain.Payment, error) {
	p := &domain.Payment{
		ID:        uuid.New(),
		BookingID: bookingID,
		UserID:    userID,
		Amount:    amount,
		Status:    domain.PaymentStatusPaid,
		Method:    method,
		PaidAt:    time.Now(),
	}
	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *PaymentUsecase) GetPayment(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *PaymentUsecase) ListUserPayments(ctx context.Context, userID uuid.UUID) ([]*domain.Payment, error) {
	return uc.repo.ListByUserID(ctx, userID)
}
