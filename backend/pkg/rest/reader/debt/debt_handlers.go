package debt

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/payvue/payvue-backend/pkg/domain/debt"
	"github.com/payvue/payvue-backend/pkg/rest/entities"
)

func (h *handler) GetAllDebts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	debts, err := h.debtService.GetAllDebts(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error_getting_debts", err.Error())
		return
	}

	response := debt.ToDebtListResponse(debts)
	respondWithJSON(w, http.StatusOK, response.Debts)
}

func (h *handler) GetDebtByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid_id", "ID must be a number")
		return
	}

	debtData, err := h.debtService.GetDebtByID(ctx, id)
	if err != nil {
		if err == debt.ErrDebtNotFound {
			respondWithError(w, http.StatusNotFound, "debt_not_found", "Deuda no encontrada")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error_getting_debt", err.Error())
		return
	}

	response := debt.ToDebtResponse(debtData)
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
