// Hand-written stub replacing protoc output. Replace with: make proto
package payment

import (
	"context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var _ = codes.OK
var _ = status.New

// ── Messages ──────────────────────────────────────────────

type Payment struct {
	Id        string  `json:"id,omitempty"`
	BookingId string  `json:"booking_id,omitempty"`
	UserId    string  `json:"user_id,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	Status    string  `json:"status,omitempty"`
	Method    string  `json:"method,omitempty"`
	PaidAt    string  `json:"paid_at,omitempty"`
}

type CreatePaymentRequest struct{ BookingId, UserId string; Amount float64; Method string }
type GetPaymentRequest struct{ PaymentId string }
type ListUserPaymentsRequest struct{ UserId string }
type GetReceiptRequest struct{ PaymentId string }
type PaymentResponse struct{ Payment *Payment }
type ListPaymentsResponse struct{ Payments []*Payment }
type ReceiptResponse struct {
	PaymentId, BookingId, UserEmail string
	Amount                          float64
	PaidAt, MovieTitle, Showtime    string
	Seats                           []string
}

// ── Server ────────────────────────────────────────────────

type PaymentServiceServer interface {
	CreatePayment(context.Context, *CreatePaymentRequest) (*PaymentResponse, error)
	GetPayment(context.Context, *GetPaymentRequest) (*PaymentResponse, error)
	ListUserPayments(context.Context, *ListUserPaymentsRequest) (*ListPaymentsResponse, error)
	GetReceipt(context.Context, *GetReceiptRequest) (*ReceiptResponse, error)
	mustEmbedUnimplementedPaymentServiceServer()
}

type UnimplementedPaymentServiceServer struct{}

func (UnimplementedPaymentServiceServer) CreatePayment(context.Context, *CreatePaymentRequest) (*PaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "CreatePayment not implemented")
}
func (UnimplementedPaymentServiceServer) GetPayment(context.Context, *GetPaymentRequest) (*PaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "GetPayment not implemented")
}
func (UnimplementedPaymentServiceServer) ListUserPayments(context.Context, *ListUserPaymentsRequest) (*ListPaymentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "ListUserPayments not implemented")
}
func (UnimplementedPaymentServiceServer) GetReceipt(context.Context, *GetReceiptRequest) (*ReceiptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "GetReceipt not implemented")
}
func (UnimplementedPaymentServiceServer) mustEmbedUnimplementedPaymentServiceServer() {}

func RegisterPaymentServiceServer(s *grpc.Server, srv PaymentServiceServer) {
	s.RegisterService(&PaymentService_ServiceDesc, srv)
}

var PaymentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "payment.PaymentService",
	HandlerType: (*PaymentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreatePayment", Handler: _PaymentService_CreatePayment_Handler},
		{MethodName: "GetPayment", Handler: _PaymentService_GetPayment_Handler},
		{MethodName: "ListUserPayments", Handler: _PaymentService_ListUserPayments_Handler},
		{MethodName: "GetReceipt", Handler: _PaymentService_GetReceipt_Handler},
	},
	Streams: []grpc.StreamDesc{},
}

func _PaymentService_CreatePayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePaymentRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(PaymentServiceServer).CreatePayment(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/payment.PaymentService/CreatePayment"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(PaymentServiceServer).CreatePayment(ctx, req.(*CreatePaymentRequest)) }
	return interceptor(ctx, in, info, handler)
}
func _PaymentService_GetPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPaymentRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(PaymentServiceServer).GetPayment(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/payment.PaymentService/GetPayment"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(PaymentServiceServer).GetPayment(ctx, req.(*GetPaymentRequest)) }
	return interceptor(ctx, in, info, handler)
}
func _PaymentService_ListUserPayments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserPaymentsRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(PaymentServiceServer).ListUserPayments(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/payment.PaymentService/ListUserPayments"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(PaymentServiceServer).ListUserPayments(ctx, req.(*ListUserPaymentsRequest)) }
	return interceptor(ctx, in, info, handler)
}
func _PaymentService_GetReceipt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReceiptRequest); if err := dec(in); err != nil { return nil, err }
	if interceptor == nil { return srv.(PaymentServiceServer).GetReceipt(ctx, in) }
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/payment.PaymentService/GetReceipt"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) { return srv.(PaymentServiceServer).GetReceipt(ctx, req.(*GetReceiptRequest)) }
	return interceptor(ctx, in, info, handler)
}

// ── Client ────────────────────────────────────────────────

type PaymentServiceClient interface {
	CreatePayment(ctx context.Context, in *CreatePaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error)
	GetPayment(ctx context.Context, in *GetPaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error)
	ListUserPayments(ctx context.Context, in *ListUserPaymentsRequest, opts ...grpc.CallOption) (*ListPaymentsResponse, error)
	GetReceipt(ctx context.Context, in *GetReceiptRequest, opts ...grpc.CallOption) (*ReceiptResponse, error)
}

type paymentServiceClient struct{ cc grpc.ClientConnInterface }

func NewPaymentServiceClient(cc grpc.ClientConnInterface) PaymentServiceClient { return &paymentServiceClient{cc} }

func (c *paymentServiceClient) CreatePayment(ctx context.Context, in *CreatePaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error) {
	out := new(PaymentResponse); err := c.cc.Invoke(ctx, "/payment.PaymentService/CreatePayment", in, out, opts...); return out, err
}
func (c *paymentServiceClient) GetPayment(ctx context.Context, in *GetPaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error) {
	out := new(PaymentResponse); err := c.cc.Invoke(ctx, "/payment.PaymentService/GetPayment", in, out, opts...); return out, err
}
func (c *paymentServiceClient) ListUserPayments(ctx context.Context, in *ListUserPaymentsRequest, opts ...grpc.CallOption) (*ListPaymentsResponse, error) {
	out := new(ListPaymentsResponse); err := c.cc.Invoke(ctx, "/payment.PaymentService/ListUserPayments", in, out, opts...); return out, err
}
func (c *paymentServiceClient) GetReceipt(ctx context.Context, in *GetReceiptRequest, opts ...grpc.CallOption) (*ReceiptResponse, error) {
	out := new(ReceiptResponse); err := c.cc.Invoke(ctx, "/payment.PaymentService/GetReceipt", in, out, opts...); return out, err
}
