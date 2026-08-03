package application

import "auth-service/internal/domain"

type AuthRepository interface {
	InsertProfile(profile *domain.Profile) error
	GetProfileByEmail(email string) (*domain.Profile, error)
	UpsertAuthToken(userID, refreshToken string) error
	GetRefreshTokenByUserID(userID string) (string, error)
}
