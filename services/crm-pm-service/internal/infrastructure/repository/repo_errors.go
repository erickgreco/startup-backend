package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrTimeout         = errors.New("query_timed_out")
	ErrInvalidData     = errors.New("invalid_data")
	ErrNotFound        = errors.New("record_not_found")
	ErrUniqueViolation = errors.New("unique_constraint_violation")
)

// ** Traduce errores de Postgres para operaciones de escritura (INSERT/UPDATE)
func mapWriteError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("repository: write_timed_out: %w", ErrTimeout)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return fmt.Errorf("repository: %s: %w", pgErr.ConstraintName, ErrUniqueViolation)
		case pgerrcode.CheckViolation, pgerrcode.ForeignKeyViolation, pgerrcode.NotNullViolation:
			return fmt.Errorf("repository: %s: %w", pgErr.ConstraintName, ErrInvalidData)
		}
	}

	return fmt.Errorf("repository: write: %w", err)
}

// ** Traduce errores de Postgres para operaciones de lectura (SELECT)
func mapReadError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("repository: read_timed_out: %w", ErrTimeout)
	}

	return fmt.Errorf("repository: read: %w", err)
}
