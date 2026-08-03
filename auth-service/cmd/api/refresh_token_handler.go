package api

import (
	"auth-service/internal/application"
	"auth-service/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
)

type RefreshTokenHandler struct {
	RefreshTokenUseCase application.RefreshTokenUseCase
}

type RefreshTokenInput struct {
	UserID       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

func NewRefreshTokenHandler(refreshTokenUseCase application.RefreshTokenUseCase) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		RefreshTokenUseCase: refreshTokenUseCase,
	}
}

func (h *RefreshTokenHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if !validateContentType(w, r) {
		return
	}

	var input RefreshTokenInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.UserID == "" || input.RefreshToken == "" {
		WriteError(w, http.StatusBadRequest, "user_id and refresh_token are required")
		return
	}

	tokenResponse, err := h.RefreshTokenUseCase.RefreshToken(input.UserID, input.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			WriteError(w, http.StatusUnauthorized, "Invalid refresh token")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to refresh token")
		return
	}

	WriteJSON(w, http.StatusOK, tokenResponse)
}
