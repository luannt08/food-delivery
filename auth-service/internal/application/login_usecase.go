package application

import (
	"auth-service/internal/domain"
	"auth-service/internal/infrastructure/token"
	"errors"
	"fmt"
)

type LoginUseCase interface {
	Login(email, password string) (domain.TokenResponse, error)
}

type loginUseCase struct {
	authRepository AuthRepository
	tokenService   token.JWTService
}

func NewLoginUseCase(authRepository AuthRepository, tokenService token.JWTService) LoginUseCase {
	return &loginUseCase{
		authRepository: authRepository,
		tokenService:   tokenService,
	}
}

func (l *loginUseCase) Login(email, password string) (domain.TokenResponse, error) {
	profile, err := l.authRepository.GetProfileByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrProfileNotFound) {
			return domain.TokenResponse{}, domain.ErrInvalidCredentials
		}
		return domain.TokenResponse{}, err
	}

	if !profile.VerifyCredentials(password) {
		fmt.Println("Invalid credentials for email:", email)
		return domain.TokenResponse{}, domain.ErrInvalidCredentials
	}

	accessToken, err := l.tokenService.GenerateToken(profile.ID)
	fmt.Println("Generated access token:", accessToken)

	if err != nil {
		return domain.TokenResponse{}, err
	}

	refreshToken, err := l.tokenService.GenerateRefreshToken(32)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	if err := l.authRepository.UpsertAuthToken(profile.ID, refreshToken); err != nil {
		return domain.TokenResponse{}, err
	}

	return domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
