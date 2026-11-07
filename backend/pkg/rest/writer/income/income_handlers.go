package income

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/payvue/payvue-backend/pkg/domain/income"
	"github.com/payvue/payvue-backend/pkg/rest/entities"
)

var validate = validator.New()

func (h *handler) CreateIncome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request entities.CreateIncomeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
		return
	}

	if err := validate.Struct(request); err != nil {
		respondWithError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	_, err := h.incomeService.CreateIncome(ctx, request.ToDomain())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error_creating_income", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, entities.MessageResponse{
		Message: "Ingreso registrado exitosamente",
	})
}

func (h *handler) UpdateIncome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid_id", "ID must be a number")
		return
	}

	var request entities.UpdateIncomeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
		return
	}

	if err := validate.Struct(request); err != nil {
		respondWithError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	_, err = h.incomeService.UpdateIncome(ctx, id, request.ToDomain())
	if err != nil {
		if err == income.ErrIncomeNotFound {
			respondWithError(w, http.StatusNotFound, "income_not_found", "Ingreso no encontrado")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error_updating_income", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, entities.MessageResponse{
		Message: "Ingreso actualizado exitosamente",
	})
}

func (h *handler) DeleteIncome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid_id", "ID must be a number")
		return
	}

	err = h.incomeService.DeleteIncome(ctx, id)
	if err != nil {
		if err == income.ErrIncomeNotFound {
			respondWithError(w, http.StatusNotFound, "income_not_found", "Ingreso no encontrado")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error_deleting_income", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, entities.MessageResponse{
		Message: "Ingreso eliminado exitosamente",
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
