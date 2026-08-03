package api

import (
	"auth-service/internal/application"
	"auth-service/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
)

type RegisterHandler struct {
	RegisterUseCase application.RegisterUseCase
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewRegisterHandler(registerUseCase application.RegisterUseCase) *RegisterHandler {
	return &RegisterHandler{
		RegisterUseCase: registerUseCase,
	}
}

func (h *RegisterHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if !validateContentType(w, r) {
		return
	}

	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validateEmailAndPassword(input.Email, input.Password); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.RegisterUseCase.Register(input.Email, input.Password); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			WriteError(w, http.StatusConflict, err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}