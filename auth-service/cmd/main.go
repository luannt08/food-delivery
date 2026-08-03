package main

import (
	"auth-service/cmd/api"
	"auth-service/internal/application"
	"auth-service/internal/infrastructure/db"
	"auth-service/internal/infrastructure/repository"
	"auth-service/internal/infrastructure/token"
	"context"
	"log"
	"os"
)

func main() {
	postgresDB := db.NewPostgresDB(context.Background())
	authRepository := repository.NewAuthRepository(postgresDB)

	jwtService := token.NewJWTService(os.Getenv("JWT_SECRET_KEY"))

	registerUseCase := application.NewRegisterUseCase(authRepository)
	loginUseCase := application.NewLoginUseCase(authRepository, jwtService)
	refreshTokenUseCase := application.NewRefreshTokenUseCase(jwtService, authRepository)

	registerHandler := api.NewRegisterHandler(registerUseCase)
	loginHandler := api.NewLoginHandler(loginUseCase)
	refreshTokenHandler := api.NewRefreshTokenHandler(refreshTokenUseCase)

	authServer := api.NewAuthServer()
	authServer.RegisterHandler("/auth/register", registerHandler.Handler)
	authServer.RegisterHandler("/auth/login", loginHandler.Handler)
	authServer.RegisterHandler("/auth/refresh", refreshTokenHandler.Handler)

	if err := authServer.Start(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
