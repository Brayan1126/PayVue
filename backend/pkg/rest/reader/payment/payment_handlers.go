package payment

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/payvue/payvue-backend/pkg/domain/payment"
	"github.com/payvue/payvue-backend/pkg/rest/entities"
	"github.com/payvue/payvue-backend/pkg/utils/fileupload"
)

func (h *handler) GetAllPayments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payments, err := h.paymentService.GetAllPayments(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error_getting_payments", err.Error())
		return
	}

	response := payment.ToPaymentListResponse(payments)
	respondWithJSON(w, http.StatusOK, response.Payments)
}

func (h *handler) GetReceipt(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")

	if filename == "" {
		respondWithError(w, http.StatusBadRequest, "invalid_filename", "Filename is required")
		return
	}

	// Construir ruta del archivo
	filePath := fileupload.GetFilePath(filename)

	// Verificar que el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		respondWithError(w, http.StatusNotFound, "file_not_found", "Receipt file not found")
		return
	}

	// Servir el archivo
	http.ServeFile(w, r, filePath)
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
