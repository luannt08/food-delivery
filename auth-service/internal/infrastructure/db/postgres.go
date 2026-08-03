package db

import (
	"auth-service/internal/domain"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context) *PostgresDB {
	pool, err := pgxpool.New(ctx, buildConfig())
	if err != nil {
		log.Fatalf("Unable to parse database config: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to Postgres")

	return &PostgresDB{
		Pool: pool,
	}
}

func buildConfig() string {
	userName := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", userName, password, host, port, dbName)
}

func (a *PostgresDB) InsertProfile(profile *domain.Profile) error {
	_, err := a.Pool.Exec(context.Background(), "INSERT INTO profile (id, email, password) VALUES ($1, $2, $3)", profile.ID, profile.Email, profile.Password)
	fmt.Println("Inserted profile:", profile, "Error:", err)
	if err != nil {
		return handlePostgresError(err)
	}
	return nil
}

func (a *PostgresDB) GetProfileByEmail(email string) (*domain.Profile, error) {
	row := a.Pool.QueryRow(context.Background(), "SELECT id, email, password FROM profile WHERE email = $1", email)
	profile := &domain.Profile{}
	err := row.Scan(&profile.ID, &profile.Email, &profile.Password)
	fmt.Println("Retrieved profile:", row, "Error:", err)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, handlePostgresError(err)
	}
	return profile, nil
}

func (a *PostgresDB) UpsertAuthToken(userID, refreshToken string) error {
	id := uuid.New().String()
	_, err := a.Pool.Exec(
		context.Background(),
		`INSERT INTO auth_token (id, user_id, refresh_token) VALUES ($1, $2, $3)
		 ON CONFLICT (user_id) DO UPDATE SET refresh_token = EXCLUDED.refresh_token, updated_at = CURRENT_TIMESTAMP`,
		id, userID, refreshToken,
	)
	if err != nil {
		return handlePostgresError(err)
	}
	return nil
}

func (a *PostgresDB) GetRefreshTokenByUserID(userID string) (string, error) {
	row := a.Pool.QueryRow(context.Background(), "SELECT refresh_token FROM auth_token WHERE user_id = $1", userID)
	var refreshToken string
	if err := row.Scan(&refreshToken); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", handlePostgresError(err)
	}
	return refreshToken, nil
}
