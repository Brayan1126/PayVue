package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/payvue/payvue-backend/cmd/app/config"
	"github.com/payvue/payvue-backend/cmd/app/container"
	writerAuth "github.com/payvue/payvue-backend/pkg/rest/writer/auth"
	writerDebt "github.com/payvue/payvue-backend/pkg/rest/writer/debt"
	writerIncome "github.com/payvue/payvue-backend/pkg/rest/writer/income"
	writerPayment "github.com/payvue/payvue-backend/pkg/rest/writer/payment"
)

func main() {
	cfg := config.Get()

	globalContainer := container.New(cfg)
	defer globalContainer.Close()

	// Crear handlers para cada módulo
	debtHandler := writerDebt.NewHandler(globalContainer.DebtService)
	incomeHandler := writerIncome.NewHandler(globalContainer.IncomeService)
	paymentHandler := writerPayment.NewHandler(globalContainer.PaymentService)
	authHandler := writerAuth.NewHandler(globalContainer.UserService)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Registrar rutas de cada módulo
	debtHandler.RouteURLs(router)
	incomeHandler.RouteURLs(router)
	paymentHandler.RouteURLs(router)
	authHandler.RouteURLs(router)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK - Writer Service"))
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting PayVue Writer API (POST/PUT/DELETE operations) on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down writer server...")
}
