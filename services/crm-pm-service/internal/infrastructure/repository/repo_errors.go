package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func mapCreateLeadDBError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("repository: create_lead_timed_out: %w", domain.ErrRepoTimeout)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return domain.ErrEmailAlreadyExists
		case pgerrcode.CheckViolation:
			return fmt.Errorf("repository: constraint_violation: %w", domain.ErrInvalidData)
		case pgerrcode.ForeignKeyViolation:
			return fmt.Errorf("repository: foreign_key_violation: %w", domain.ErrInvalidData)
		case pgerrcode.NotNullViolation:
			return fmt.Errorf("repository: missing_required_field: %w", domain.ErrInvalidData)
		}
	}

	return fmt.Errorf("repository: create_lead: %w", err)
}
