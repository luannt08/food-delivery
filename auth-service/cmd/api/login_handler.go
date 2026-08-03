package api

import (
	"auth-service/internal/application"
	"encoding/json"
	"net/http"
)

type LoginHandler struct {
	LoginUseCase application.LoginUseCase
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLoginHandler(loginUseCase application.LoginUseCase) *LoginHandler {
	return &LoginHandler{
		LoginUseCase: loginUseCase,
	}
}

func (h *LoginHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if !validateContentType(w, r) {
		return
	}

	var input LoginInput
	er := json.NewDecoder(r.Body).Decode(&input)

	if er != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tokenResponse, err := h.LoginUseCase.Login(input.Email, input.Password)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, tokenResponse)
}
