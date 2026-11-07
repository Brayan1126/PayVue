package income

import (
	"github.com/payvue/payvue-backend/pkg/domain/income"
	"github.com/payvue/payvue-backend/pkg/rest"
)

type handler struct {
	incomeService income.Service
}

func NewHandler(incomeService income.Service) rest.Handler {
	return &handler{
		incomeService: incomeService,
	}
}
