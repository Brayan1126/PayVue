package debt

import (
	"github.com/payvue/payvue-backend/pkg/domain/debt"
	"github.com/payvue/payvue-backend/pkg/rest"
)

type handler struct {
	debtService debt.Service
}

func NewHandler(debtService debt.Service) rest.Handler {
	return &handler{
		debtService: debtService,
	}
}
