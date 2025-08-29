package paymentHandler

// import (
// 	"context"
// 	"microService/modules/payment/domain"
// 	"microService/modules/payment/paymentPb"
// 	"microService/modules/payment/paymentUsecase"
// )

// type (
// 	PaymentGrpcHandler struct {
// 		paymentUsecase paymentUsecase.PaymentUsecase
// 		paymentPb.UnimplementedPaymentServiceServer
// 	}
// )

// func NewpaymentGrpcHandler(paymentUsecase paymentUsecase.PaymentUsecase) *PaymentGrpcHandler {
// 	return &PaymentGrpcHandler{paymentUsecase: paymentUsecase}
// }

// func (s *PaymentGrpcHandler) CreatePayment(ctx context.Context, in *paymentPb.CreatePaymentRequest) (*paymentPb.CreatePaymentResponse, error) {
// 	p, err := s.paymentUsecase.CreatePendingPayment(ctx, in.GetOrderId(), in.GetUserId(), in.GetAmount(), in.GetCurrency())
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &paymentPb.CreatePaymentResponse{
// 		PaymentId: p.ID,
// 		Status:    p.Status,
// 		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
// 	}, nil
// }

// func (s *PaymentGrpcHandler) GetPayment(ctx context.Context, in *paymentPb.GetPaymentRequest) (*paymentPb.GetPaymentResponse, error) {
// 	return &paymentPb.GetPaymentResponse{PaymentId: in.GetPaymentId(), Status: domain.PaymentStatusPending}, nil
// }
