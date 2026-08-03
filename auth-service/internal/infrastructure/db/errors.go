package db

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUniqueViolation = errors.New("unique constraint violation")
	ErrNotFound        = errors.New("not found")
)

func handlePostgresError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return ErrUniqueViolation
		}
	}
	
	return fmt.Errorf("postgres error: %v", err)
}