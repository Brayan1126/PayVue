package debt

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/payvue/payvue-backend/pkg/domain/debt"
	"github.com/payvue/payvue-backend/pkg/rest/entities"
)

var validate = validator.New()

func (h *handler) CreateDebt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request entities.CreateDebtRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
		return
	}

	if err := validate.Struct(request); err != nil {
		respondWithError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	_, err := h.debtService.CreateDebt(ctx, request.ToDomain())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error_creating_debt", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, entities.MessageResponse{
		Message: "Deuda registrada exitosamente",
	})
}

func (h *handler) UpdateDebt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid_id", "ID must be a number")
		return
	}

	var request entities.UpdateDebtRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
		return
	}

	if err := validate.Struct(request); err != nil {
		respondWithError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	_, err = h.debtService.UpdateDebt(ctx, id, request.ToDomain())
	if err != nil {
		if err == debt.ErrDebtNotFound {
			respondWithError(w, http.StatusNotFound, "debt_not_found", "Deuda no encontrada")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error_updating_debt", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, entities.MessageResponse{
		Message: "Deuda actualizada exitosamente",
	})
}

func (h *handler) DeleteDebt(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid_id", "ID must be a number")
		return
	}

	err = h.debtService.DeleteDebt(ctx, id)
	if err != nil {
		if err == debt.ErrDebtNotFound {
			respondWithError(w, http.StatusNotFound, "debt_not_found", "Deuda no encontrada")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error_deleting_debt", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, entities.MessageResponse{
		Message: "Deuda eliminada exitosamente",
	})
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
