package container

import (
	"database/sql"
	"log"

	"github.com/payvue/payvue-backend/cmd/app/config"
	"github.com/payvue/payvue-backend/pkg/domain/debt"
	"github.com/payvue/payvue-backend/pkg/domain/income"
	"github.com/payvue/payvue-backend/pkg/domain/payment"
	"github.com/payvue/payvue-backend/pkg/domain/user"
	"github.com/payvue/payvue-backend/pkg/repository/database"
	debtRepo "github.com/payvue/payvue-backend/pkg/repository/debt"
	incomeRepo "github.com/payvue/payvue-backend/pkg/repository/income"
	paymentRepo "github.com/payvue/payvue-backend/pkg/repository/payment"
	userRepo "github.com/payvue/payvue-backend/pkg/repository/user"
)

type Container struct {
	DebtService    debt.Service
	IncomeService  income.Service
	PaymentService payment.Service
	UserService    user.Service
	DB             *sql.DB
}

func New(cfg config.Config) *Container {
	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Debt
	debtRepository := debtRepo.NewRepository(db)
	debtContainer := &debt.Container{
		Repository: debtRepository,
	}
	debtService := debt.New(debtContainer)

	// Income
	incomeRepository := incomeRepo.NewRepository(db)
	incomeContainer := &income.Container{
		Repository: incomeRepository,
	}
	incomeService := income.New(incomeContainer)

	// Payment
	paymentRepository := paymentRepo.NewRepository(db)
	paymentContainer := &payment.Container{
		Repository: paymentRepository,
	}
	paymentService := payment.New(paymentContainer)

	// User
	userRepository := userRepo.NewRepository(db)
	userContainer := &user.Container{
		Repository: userRepository,
	}
	userService := user.New(userContainer)

	return &Container{
		DebtService:    debtService,
		IncomeService:  incomeService,
		PaymentService: paymentService,
		UserService:    userService,
		DB:             db,
	}
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
