package application

import (
	"auth-service/internal/domain"

	"github.com/google/uuid"
)

type RegisterUseCase interface {
	Register(email, password string) error
}

type registerUseCase struct {
	authRepository AuthRepository
}

func NewRegisterUseCase(authRepository AuthRepository) RegisterUseCase {
	return &registerUseCase{
		authRepository: authRepository,
	}
}

func (r *registerUseCase) Register(email, password string) error {
	hashedPassword, err := domain.HashPassword(password)
	if err != nil {
		return err
	}

	profile := domain.NewProfile(uuid.New().String(), email, hashedPassword)

	if err := r.authRepository.InsertProfile(profile); err != nil {
		return err
	}

	return nil
}
