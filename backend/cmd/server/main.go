package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/payvue/payvue-backend/cmd/app/config"
	"github.com/payvue/payvue-backend/cmd/app/container"
	"github.com/payvue/payvue-backend/pkg/domain/debt"
	"github.com/payvue/payvue-backend/pkg/domain/income"
	"github.com/payvue/payvue-backend/pkg/domain/payment"
	"github.com/payvue/payvue-backend/pkg/domain/user"
	"github.com/payvue/payvue-backend/pkg/rest/entities"
	"github.com/payvue/payvue-backend/pkg/utils/fileupload"
)

var validate = validator.New()

func main() {
	cfg := config.Get()

	globalContainer := container.New(cfg)
	defer globalContainer.Close()

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-User-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Auth routes
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", makeAuthRegisterHandler(globalContainer.UserService))
		r.Post("/login", makeAuthLoginHandler(globalContainer.UserService))
		r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
			respondWithJSON(w, http.StatusOK, map[string]string{"message": "Sesión cerrada exitosamente"})
		})
	})

	// Income routes (combined reader + writer)
	router.Route("/finances/income", func(r chi.Router) {
		r.Get("/", makeGetAllIncomesHandler(globalContainer.IncomeService))
		r.Get("/{id}", makeGetIncomeByIDHandler(globalContainer.IncomeService))
		r.Post("/", makeCreateIncomeHandler(globalContainer.IncomeService))
		r.Put("/{id}", makeUpdateIncomeHandler(globalContainer.IncomeService))
		r.Delete("/{id}", makeDeleteIncomeHandler(globalContainer.IncomeService))
	})

	// Debt routes (combined reader + writer)
	router.Route("/finances/debt", func(r chi.Router) {
		r.Get("/", makeGetAllDebtsHandler(globalContainer.DebtService))
		r.Get("/{id}", makeGetDebtByIDHandler(globalContainer.DebtService))
		r.Post("/", makeCreateDebtHandler(globalContainer.DebtService))
		r.Put("/{id}", makeUpdateDebtHandler(globalContainer.DebtService))
		r.Delete("/{id}", makeDeleteDebtHandler(globalContainer.DebtService))
	})

	// Payment routes (combined reader + writer)
	router.Route("/finances/payment", func(r chi.Router) {
		r.Get("/", makeGetAllPaymentsHandler(globalContainer.PaymentService))
		r.Get("/receipt/{filename}", makeGetReceiptHandler())
		r.Post("/", makeCreatePaymentHandler(globalContainer.PaymentService))
		r.Delete("/{id}", makeDeletePaymentHandler(globalContainer.PaymentService))
	})

	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK - PayVue API Server"))
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Starting PayVue Unified API Server on port %s", cfg.Port)
		log.Printf("📍 Environment: %s", cfg.Environment)
		log.Println("📋 Available endpoints:")
		log.Println("   - POST /auth/register, /auth/login, /auth/logout")
		log.Println("   - GET/POST/PUT/DELETE /finances/income/*")
		log.Println("   - GET/POST/PUT/DELETE /finances/debt/*")
		log.Println("   - GET/POST/DELETE /finances/payment/*")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

// Helper functions
func respondWithError(w http.ResponseWriter, code int, error string, message string) {
	respondWithJSON(w, code, map[string]string{"error": error, "message": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getUserIDFromHeader(r *http.Request) int {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		// Also check query param as fallback
		userIDStr = r.URL.Query().Get("user_id")
	}
	userID, _ := strconv.Atoi(userIDStr)
	return userID
}

// Auth handlers
func makeAuthRegisterHandler(userService user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request entities.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
			return
		}
		if err := validate.Struct(request); err != nil {
			respondWithError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}
		u, err := userService.Register(r.Context(), request.ToDomain())
		if err != nil {
			if err == user.ErrEmailAlreadyExists {
				respondWithError(w, http.StatusBadRequest, "email_already_exists", "El correo ya está registrado")
				return
			}
			respondWithError(w, http.StatusInternalServerError, "error_registering_user", err.Error())
			return
		}
		respondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "Usuario registrado exitosamente",
			"user_id": u.ID,
			"email":   u.Email,
		})
	}
}

func makeAuthLoginHandler(userService user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request entities.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
			return
		}
		if err := validate.Struct(request); err != nil {
			respondWithError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}
		u, err := userService.Login(r.Context(), request.ToDomain())
		if err != nil {
			if err == user.ErrInvalidCredentials {
				respondWithError(w, http.StatusUnauthorized, "invalid_credentials", "Credenciales inválidas")
				return
			}
			respondWithError(w, http.StatusInternalServerError, "error_login", err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Inicio de sesión exitoso",
			"user_id": u.ID,
			"email":   u.Email,
		})
	}
}

// Income handlers
func makeGetAllIncomesHandler(incomeService income.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromHeader(r)
		var incomes []income.Income
		var err error

		if userID > 0 {
			incomes, err = incomeService.GetIncomesByUserID(r.Context(), userID)
		} else {
			incomes, err = incomeService.GetAllIncomes(r.Context())
		}

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error_getting_incomes", err.Error())
			return
		}
		responses := make([]income.IncomeResponse, len(incomes))
		for i, inc := range incomes {
			responses[i] = income.ToIncomeResponse(&inc)
		}
		respondWithJSON(w, http.StatusOK, responses)
	}
}

func makeGetIncomeByIDHandler(incomeService income.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		inc, err := incomeService.GetIncomeByID(r.Context(), id)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "income_not_found", "Ingreso no encontrado")
			return
		}
		respondWithJSON(w, http.StatusOK, income.ToIncomeResponse(inc))
	}
}

func makeCreateIncomeHandler(incomeService income.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromHeader(r)

		var request entities.CreateIncomeRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
			return
		}
		if err := validate.Struct(request); err != nil {
			respondWithError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}

		domainReq := request.ToDomain()
		domainReq.UserID = userID

		inc, err := incomeService.CreateIncome(r.Context(), domainReq)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error_creating_income", err.Error())
			return
		}
		respondWithJSON(w, http.StatusCreated, income.ToIncomeResponse(inc))
	}
}

func makeUpdateIncomeHandler(incomeService income.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		var request entities.UpdateIncomeRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
			return
		}
		inc, err := incomeService.UpdateIncome(r.Context(), id, request.ToDomain())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error_updating_income", err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, income.ToIncomeResponse(inc))
	}
}

func makeDeleteIncomeHandler(incomeService income.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		if err := incomeService.DeleteIncome(r.Context(), id); err != nil {
			respondWithError(w, http.StatusInternalServerError, "error_deleting_income", err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Ingreso eliminado exitosamente"})
	}
}

// Debt handlers
func makeGetAllDebtsHandler(debtService debt.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromHeader(r)
		var debts []debt.Debt
		var err error

		if userID > 0 {
			debts, err = debtService.GetDebtsByUserID(r.Context(), userID)
		} else {
			debts, err = debtService.GetAllDebts(r.Context())
		}

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error_getting_debts", err.Error())
			return
		}
		responses := make([]debt.DebtResponse, len(debts))
		for i, d := range debts {
			responses[i] = debt.ToDebtResponse(&d)
		}
		respondWithJSON(w, http.StatusOK, responses)
	}
}

func makeGetDebtByIDHandler(debtService debt.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		d, err := debtService.GetDebtByID(r.Context(), id)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "debt_not_found", "Deuda no encontrada")
			return
		}
		respondWithJSON(w, http.StatusOK, debt.ToDebtResponse(d))
	}
}

func makeCreateDebtHandler(debtService debt.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromHeader(r)

		var request entities.CreateDebtRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
			return
		}
		if err := validate.Struct(request); err != nil {
			respondWithError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}

		domainReq := request.ToDomain()
		domainReq.UserID = userID

		d, err := debtService.CreateDebt(r.Context(), domainReq)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error_creating_debt", err.Error())
			return
		}
		respondWithJSON(w, http.StatusCreated, debt.ToDebtResponse(d))
	}
}

func makeUpdateDebtHandler(debtService debt.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		var request entities.UpdateDebtRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
			return
		}
		d, err := debtService.UpdateDebt(r.Context(), id, request.ToDomain())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error_updating_debt", err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, debt.ToDebtResponse(d))
	}
}

func makeDeleteDebtHandler(debtService debt.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		if err := debtService.DeleteDebt(r.Context(), id); err != nil {
			respondWithError(w, http.StatusInternalServerError, "error_deleting_debt", err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Deuda eliminada exitosamente"})
	}
}

// Payment handlers
func makeGetAllPaymentsHandler(paymentService payment.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromHeader(r)
		var payments []payment.PaymentWithDebt
		var err error

		if userID > 0 {
			payments, err = paymentService.GetPaymentsByUserID(r.Context(), userID)
		} else {
			payments, err = paymentService.GetAllPayments(r.Context())
		}

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error_getting_payments", err.Error())
			return
		}
		responses := make([]payment.PaymentResponse, len(payments))
		for i, p := range payments {
			responses[i] = payment.ToPaymentResponse(p)
		}
		respondWithJSON(w, http.StatusOK, responses)
	}
}

func makeGetReceiptHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := chi.URLParam(r, "filename")
		filePath := fileupload.GetFilePath(filename)
		http.ServeFile(w, r, filePath)
	}
}

func makeCreatePaymentHandler(paymentService payment.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromHeader(r)
		r.ParseMultipartForm(fileupload.MaxFileSize)

		amountStr := r.FormValue("amount")
		amount, _ := strconv.ParseFloat(amountStr, 64)
		debtIDStr := r.FormValue("debt_id")
		debtID, _ := strconv.Atoi(debtIDStr)
		date := r.FormValue("date")

		var filename string
		file, header, err := r.FormFile("receipt")
		if err == nil {
			defer file.Close()
			filename, _ = fileupload.SaveFile(file, header)
		}

		request := payment.CreatePaymentRequest{
			UserID: userID,
			Amount: amount,
			DebtID: debtID,
			Date:   date,
		}

		p, err := paymentService.CreatePayment(r.Context(), request, filename)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error_creating_payment", err.Error())
			return
		}
		respondWithJSON(w, http.StatusCreated, p)
	}
}

func makeDeletePaymentHandler(paymentService payment.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		if err := paymentService.DeletePayment(r.Context(), id); err != nil {
			respondWithError(w, http.StatusInternalServerError, "error_deleting_payment", err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Pago eliminado exitosamente"})
	}
}
