package income

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/payvue/payvue-backend/pkg/domain/income"
	"github.com/payvue/payvue-backend/pkg/rest/entities"
)

func (h *handler) GetAllIncomes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	incomes, err := h.incomeService.GetAllIncomes(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error_getting_incomes", err.Error())
		return
	}

	response := income.ToIncomeListResponse(incomes)
	respondWithJSON(w, http.StatusOK, response.Incomes)
}

func (h *handler) GetIncomeByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid_id", "ID must be a number")
		return
	}

	incomeData, err := h.incomeService.GetIncomeByID(ctx, id)
	if err != nil {
		if err == income.ErrIncomeNotFound {
			respondWithError(w, http.StatusNotFound, "income_not_found", "Ingreso no encontrado")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error_getting_income", err.Error())
		return
	}

	response := income.ToIncomeResponse(incomeData)
	respondWithJSON(w, http.StatusOK, response)
}

func respondWithError(w http.ResponseWriter, code int, error string, message string) {
	respondWithJSON(w, code, entities.ErrorResponse{
		Error:   error,
		Message: message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
