package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// TODO: hay errores que se generan pero no estan siendo mapeados, pgerrcode tiene helpers como pgerrcode.IsDataException()
var (
	ErrTimeout         = errors.New("query_timed_out")
	ErrInvalidData     = errors.New("invalid_data")
	ErrUniqueViolation = errors.New("unique_constraint_violation")
)

// * Traduce errores de Postgres para operaciones de escritura (INSERT/UPDATE)
func mapWriteError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("repository: write_timed_out: %w", ErrTimeout)
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return fmt.Errorf("repository: %s: %w", pgErr.ConstraintName, ErrUniqueViolation)
		case pgerrcode.CheckViolation, pgerrcode.ForeignKeyViolation, pgerrcode.NotNullViolation:
			return fmt.Errorf("repository: %s: %w", pgErr.ConstraintName, ErrInvalidData)
		}
	}

	return fmt.Errorf("repository: write: %w", err)
}

// * Traduce el error generico de infraestructura a un error de dominio
func translateWriteError(err error) error {
	mappedErr := mapWriteError(err)

	if errors.Is(mappedErr, ErrUniqueViolation) {
		return domain.ErrAlreadyExists
	}
	if errors.Is(mappedErr, ErrInvalidData) {
		return fmt.Errorf("%w: %w", domain.ErrInvalidData, mappedErr)
	}
	if errors.Is(mappedErr, ErrTimeout) {
		return domain.ErrRepoTimeout
	}

	return fmt.Errorf("%w: %w", domain.ErrInternal, mappedErr)
}

// * Traduce errores de Postgres para operaciones de lectura (SELECT)
func mapReadError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("repository: read_timed_out: %w", ErrTimeout)
	}

	return fmt.Errorf("repository: read: %w", err)
}

// * Traduce el error generico de infraestructura a un error de dominio
func translateReadError(err error) error {
	mappedErr := mapReadError(err)

	if errors.Is(mappedErr, domain.ErrNotFound) {
		return domain.ErrNotFound
	}
	if errors.Is(mappedErr, ErrTimeout) {
		return domain.ErrRepoTimeout
	}

	return fmt.Errorf("%w: %w", domain.ErrInternal, mappedErr)
}

// * UPDATE -> RETURNING sin filas coincidientes -> el scan de :one devuelve pgx.ErrNoRows
func translateUpdateError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}

	mappedErr := mapWriteError(err)

	if errors.Is(mappedErr, ErrUniqueViolation) {
		return domain.ErrAlreadyExists
	}
	if errors.Is(mappedErr, ErrInvalidData) {
		return domain.ErrInvalidData
	}
	if errors.Is(mappedErr, ErrTimeout) {
		return domain.ErrRepoTimeout
	}

	return fmt.Errorf("%w: %w", domain.ErrInternal, mappedErr)
}
