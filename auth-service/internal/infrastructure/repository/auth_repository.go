package repository

import (
	"auth-service/internal/application"
	"auth-service/internal/domain"
	"auth-service/internal/infrastructure/db"
	"errors"
)

type authRepository struct {
	PostgresDB *db.PostgresDB
}

func NewAuthRepository(postgresDB *db.PostgresDB) application.AuthRepository {
	return &authRepository{
		PostgresDB: postgresDB,
	}
}

func (r *authRepository) GetProfileByEmail(email string) (*domain.Profile, error) {
	profile, err := r.PostgresDB.GetProfileByEmail(email)
	if errors.Is(err, db.ErrNotFound) {
		return nil, domain.ErrProfileNotFound
	}
	return profile, err
}

func (r *authRepository) InsertProfile(profile *domain.Profile) error {
	err := r.PostgresDB.InsertProfile(profile)
	if errors.Is(err, db.ErrUniqueViolation) {
		return domain.ErrEmailAlreadyExists
	}
	return err
}

func (r *authRepository) UpsertAuthToken(userID, refreshToken string) error {
	return r.PostgresDB.UpsertAuthToken(userID, refreshToken)
}

func (r *authRepository) GetRefreshTokenByUserID(userID string) (string, error) {
	refreshToken, err := r.PostgresDB.GetRefreshTokenByUserID(userID)
	if errors.Is(err, db.ErrNotFound) {
		return "", domain.ErrInvalidToken
	}
	return refreshToken, err
}