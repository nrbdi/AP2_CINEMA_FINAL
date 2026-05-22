package grpc

import (
	"context"

	"cinema/payment-service/internal/domain"
	"cinema/payment-service/internal/usecase"
	pb "cinema/proto/payment"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUsecase
}

func NewPaymentHandler(uc *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{uc: uc}
}

func mapErr(err error) error {
	if err == domain.ErrNotFound {
		return status.Error(codes.NotFound, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}

func paymentToPB(p *domain.Payment) *pb.Payment {
	return &pb.Payment{
		Id:        p.ID.String(),
		BookingId: p.BookingID.String(),
		UserId:    p.UserID.String(),
		Amount:    p.Amount,
		Status:    string(p.Status),
		Method:    p.Method,
		PaidAt:    p.PaidAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (h *PaymentHandler) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.PaymentResponse, error) {
	bookingID, _ := uuid.Parse(req.BookingId)
	userID, _ := uuid.Parse(req.UserId)
	p, err := h.uc.CreatePayment(ctx, bookingID, userID, req.Amount, req.Method)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.PaymentResponse{Payment: paymentToPB(p)}, nil
}

func (h *PaymentHandler) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.PaymentResponse, error) {
	id, _ := uuid.Parse(req.PaymentId)
	p, err := h.uc.GetPayment(ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.PaymentResponse{Payment: paymentToPB(p)}, nil
}

func (h *PaymentHandler) ListUserPayments(ctx context.Context, req *pb.ListUserPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	userID, _ := uuid.Parse(req.UserId)
	payments, err := h.uc.ListUserPayments(ctx, userID)
	if err != nil {
		return nil, mapErr(err)
	}
	pbPayments := make([]*pb.Payment, len(payments))
	for i, p := range payments {
		pbPayments[i] = paymentToPB(p)
	}
	return &pb.ListPaymentsResponse{Payments: pbPayments}, nil
}

func (h *PaymentHandler) GetReceipt(ctx context.Context, req *pb.GetReceiptRequest) (*pb.ReceiptResponse, error) {
	id, _ := uuid.Parse(req.PaymentId)
	p, err := h.uc.GetPayment(ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	return &pb.ReceiptResponse{
		PaymentId: p.ID.String(),
		BookingId: p.BookingID.String(),
		Amount:    p.Amount,
		PaidAt:    p.PaidAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}
