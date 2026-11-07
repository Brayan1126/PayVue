package payment

import (
	"github.com/payvue/payvue-backend/pkg/domain/payment"
	"github.com/payvue/payvue-backend/pkg/rest"
)

type handler struct {
	paymentService payment.Service
}

func NewHandler(paymentService payment.Service) rest.Handler {
	return &handler{
		paymentService: paymentService,
	}
}
