package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/payvue/payvue-backend/pkg/domain/user"
	"github.com/payvue/payvue-backend/pkg/rest/entities"
)

var validate = validator.New()

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request entities.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
		return
	}

	if err := validate.Struct(request); err != nil {
		respondWithError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	_, err := h.userService.Register(ctx, request.ToDomain())
	if err != nil {
		if err == user.ErrEmailAlreadyExists {
			respondWithError(w, http.StatusBadRequest, "email_already_exists", "El correo ya está registrado")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error_registering_user", err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, entities.MessageResponse{
		Message: "Usuario registrado exitosamente",
	})
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request entities.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "error_decoding_json", err.Error())
		return
	}

	if err := validate.Struct(request); err != nil {
		respondWithError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	_, err := h.userService.Login(ctx, request.ToDomain())
	if err != nil {
		if err == user.ErrInvalidCredentials {
			respondWithError(w, http.StatusUnauthorized, "invalid_credentials", "Credenciales inválidas")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "error_login", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, entities.MessageResponse{
		Message: "Inicio de sesión exitoso",
	})
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, entities.MessageResponse{
		Message: "Sesión cerrada exitosamente",
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
