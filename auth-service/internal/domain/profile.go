package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Profile struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

func NewProfile(id, email, hashedPassword string) *Profile {
	return &Profile{
		ID:       id,
		Email:    email,
		Password: hashedPassword,
	}
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (p *Profile) VerifyCredentials(password string) bool {
	fmt.Printf("Verifying credentials for original pass : %s - checking pass: %s\n", p.Password, password)
	return bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password)) == nil
}
