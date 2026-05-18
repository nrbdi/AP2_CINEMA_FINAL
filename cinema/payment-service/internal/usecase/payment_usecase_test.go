package usecase_test

import (
	"context"
	"testing"
	"time"

	"cinema/payment-service/internal/domain"
	"cinema/payment-service/internal/usecase"
	"github.com/google/uuid"
)

// ── Mock ─────────────────────────────────────────────────
type mockPaymentRepo struct {
	payments map[uuid.UUID]*domain.Payment
}

func newMockRepo() *mockPaymentRepo {
	return &mockPaymentRepo{payments: make(map[uuid.UUID]*domain.Payment)}
}

func (m *mockPaymentRepo) Create(_ context.Context, p *domain.Payment) error {
	m.payments[p.ID] = p
	return nil
}
func (m *mockPaymentRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Payment, error) {
	p, ok := m.payments[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return p, nil
}
func (m *mockPaymentRepo) GetByBookingID(_ context.Context, bookingID uuid.UUID) (*domain.Payment, error) {
	for _, p := range m.payments {
		if p.BookingID == bookingID {
			return p, nil
		}
	}
	return nil, domain.ErrNotFound
}
func (m *mockPaymentRepo) ListByUserID(_ context.Context, userID uuid.UUID) ([]*domain.Payment, error) {
	var result []*domain.Payment
	for _, p := range m.payments {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}
func (m *mockPaymentRepo) UpdateStatus(_ context.Context, id uuid.UUID, status domain.PaymentStatus) error {
	if p, ok := m.payments[id]; ok {
		p.Status = status
	}
	return nil
}

// ── Tests ─────────────────────────────────────────────────
func TestCreatePayment_Success(t *testing.T) {
	repo := newMockRepo()
	uc := usecase.NewPaymentUsecase(repo)

	bookingID := uuid.New()
	userID := uuid.New()

	p, err := uc.CreatePayment(context.Background(), bookingID, userID, 25.00, "card")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.ID == uuid.Nil {
		t.Error("expected non-nil payment ID")
	}
	if p.Amount != 25.00 {
		t.Errorf("expected amount 25.00, got %.2f", p.Amount)
	}
	if p.Status != domain.PaymentStatusPaid {
		t.Errorf("expected status 'paid', got %s", p.Status)
	}
	if p.PaidAt.IsZero() {
		t.Error("expected non-zero paid_at")
	}
	if _, ok := repo.payments[p.ID]; !ok {
		t.Error("payment not stored in repo")
	}
}

func TestGetPayment_NotFound(t *testing.T) {
	uc := usecase.NewPaymentUsecase(newMockRepo())
	_, err := uc.GetPayment(context.Background(), uuid.New())
	if err != domain.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestListUserPayments_Empty(t *testing.T) {
	uc := usecase.NewPaymentUsecase(newMockRepo())
	payments, err := uc.ListUserPayments(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 0 {
		t.Errorf("expected 0 payments, got %d", len(payments))
	}
}

func TestListUserPayments_Multiple(t *testing.T) {
	repo := newMockRepo()
	uc := usecase.NewPaymentUsecase(repo)
	userID := uuid.New()
	otherID := uuid.New()

	uc.CreatePayment(context.Background(), uuid.New(), userID, 10.0, "card")
	uc.CreatePayment(context.Background(), uuid.New(), userID, 20.0, "card")
	uc.CreatePayment(context.Background(), uuid.New(), otherID, 30.0, "card") // different user

	payments, err := uc.ListUserPayments(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payments) != 2 {
		t.Errorf("expected 2 payments for user, got %d", len(payments))
	}
}

func TestPayment_PaidAtIsRecent(t *testing.T) {
	uc := usecase.NewPaymentUsecase(newMockRepo())
	before := time.Now().Add(-time.Second)
	p, _ := uc.CreatePayment(context.Background(), uuid.New(), uuid.New(), 15.0, "system")
	after := time.Now().Add(time.Second)

	if p.PaidAt.Before(before) || p.PaidAt.After(after) {
		t.Errorf("paid_at %v not in expected range [%v, %v]", p.PaidAt, before, after)
	}
}
