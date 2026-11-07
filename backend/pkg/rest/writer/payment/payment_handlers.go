package payment

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/payvue/payvue-backend/pkg/domain/payment"
	"github.com/payvue/payvue-backend/pkg/rest/entities"
	"github.com/payvue/payvue-backend/pkg/utils/fileupload"
)

func (h *handler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse multipart form (10 MB max)
	err := r.ParseMultipartForm(fileupload.MaxFileSize)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error_parsing_form", err.Error())
		return
	}

	// Obtener campos del formulario
	amountStr := r.FormValue("amount")
	debtIDStr := r.FormValue("debt_id")
	date := r.FormValue("date")

	// Validar campos requeridos
	if amountStr == "" || debtIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing_fields", "Amount and debt_id are required")
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid_amount", "Amount must be a number")
		return
	}

	debtID, err := strconv.Atoi(debtIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid_debt_id", "Debt ID must be a number")
		return
	}

	// Obtener archivo
	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "file_required", "Receipt file is required")
		return
	}
	defer file.Close()

	// Guardar archivo
	filename, err := fileupload.SaveFile(file, header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error_saving_file", err.Error())
		return
	}

	// Crear request
	request := payment.CreatePaymentRequest{
		Amount: amount,
		DebtID: debtID,
		Date:   date,
	}

	// Crear pago
	_, err = h.paymentService.CreatePayment(ctx, request, filename)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error_creating_payment", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, entities.MessageResponse{
		Message: "Pago registrado exitosamente",
	})
}

func (h *handler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid_id", "ID must be a number")
		return
	}

	err = h.paymentService.DeletePayment(ctx, id)
	if err != nil {
		if err == payment.ErrPaymentNotFound {
			respondWithError(w, http.StatusNotFound, "payment_not_found", "Pago no encontrado")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error_deleting_payment", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, entities.MessageResponse{
		Message: "Pago eliminado exitosamente",
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
