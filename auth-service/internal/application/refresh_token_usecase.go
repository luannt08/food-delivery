package application

import (
	"auth-service/internal/domain"
	jwt "auth-service/internal/infrastructure/token"
)

type RefreshTokenUseCase interface {
	RefreshToken(userID, refreshToken string) (domain.TokenResponse, error)
}

type refreshTokenUseCase struct {
	JwtService     jwt.JWTService
	AuthRepository AuthRepository
}

func NewRefreshTokenUseCase(jwtService jwt.JWTService, authRepository AuthRepository) RefreshTokenUseCase {
	return &refreshTokenUseCase{
		JwtService:     jwtService,
		AuthRepository: authRepository,
	}
}

func (r *refreshTokenUseCase) RefreshToken(userID, refreshToken string) (domain.TokenResponse, error) {
	storedRefreshToken, err := r.AuthRepository.GetRefreshTokenByUserID(userID)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	if storedRefreshToken != refreshToken {
		return domain.TokenResponse{}, domain.ErrInvalidToken
	}

	newToken, err := r.JwtService.GenerateToken(userID)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	newRefreshToken, err := r.JwtService.GenerateRefreshToken(32)
	if err != nil {
		return domain.TokenResponse{}, err
	}

	if err := r.AuthRepository.UpsertAuthToken(userID, newRefreshToken); err != nil {
		return domain.TokenResponse{}, err
	}

	return domain.TokenResponse{AccessToken: newToken, RefreshToken: newRefreshToken}, nil
}
