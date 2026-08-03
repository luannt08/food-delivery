package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID string) (string, error)
	GenerateRefreshToken(length int) (string, error)
	ValidateToken(token string) (string, error)
}

const (
	Issuer         = "food-delivery"
	Audience       = "restaurant-service"
	AccessTokenTTL = 15 * time.Minute
)

var ErrInvalidToken = errors.New("invalid token")

type jwtService struct {
	secretKey string
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey: secretKey,
	}
}

func (j *jwtService) GenerateToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    Issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userID,
		Audience:  []string{Audience},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.secretKey))
	return token, err
}

func (j *jwtService) GenerateRefreshToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (j *jwtService) ValidateToken(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("unexpected signing method: %v", t.Header["alg"])
			return nil, errors.New("Internal server error")
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if !parsed.Valid {
		return "", ErrInvalidToken
	}

	return claims.Subject, nil
}
