package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	queryTimeDuration = time.Second * 5
)

type LeadRepository struct {
	pool *pgxpool.Pool
}

func NewLeadRepository(pool *pgxpool.Pool) *LeadRepository {
	return &LeadRepository{
		pool: pool,
	}
}

func (r *LeadRepository) CreateLead(ctx context.Context, lead *domain.LeadsModel) error {
	query := `
		INSERT INTO leads (
			id, name, email, message, status, source
		)
		VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	err := r.pool.QueryRow(
		ctx,
		query,
		lead.ID,
		lead.Name,
		lead.Email,
		lead.Message,
		lead.Status,
		lead.Source,
	).Scan(
		&lead.CreatedAt,
		&lead.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("leads_repo: create_lead: %w", translateWriteError(err))
	}

	return nil
}

// ** Es especifico de leads: traduce el error generico de infraestructura a un error
// ** de dominio de leads
func translateWriteError(err error) error {
	mappedErr := mapWriteError(err)

	if errors.Is(mappedErr, ErrUniqueViolation) {
		return domain.ErrEmailAlreadyExists
	}
	if errors.Is(mappedErr, ErrInvalidData) {
		return domain.ErrInvalidData
	}
	if errors.Is(mappedErr, ErrTimeout) {
		return domain.ErrRepoTimeout
	}

	return mappedErr
}

func (r *LeadRepository) GetLeadByStatus(ctx context.Context, status domain.LeadStatus) (*domain.LeadsModel, error) {
	var lead domain.LeadsModel

	query := `
		SELECT id, name, email, message, status, source, created_at, updated_at
		FROM leads
		WHERE created_at = $1
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	err := r.pool.QueryRow(
		ctx,
		query,
		status,
	).Scan(
		&lead.ID,
		&lead.Name,
		&lead.Email,
		&lead.Message,
		&lead.Status,
		&lead.Source,
		&lead.CreatedAt,
		&lead.UpdatedAt,
	)
	if err != nil {
		mappedErr := mapReadError(err)
		if errors.Is(mappedErr, ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("leads_repo: get_lead_by_date: %w", mappedErr)
	}

	return &lead, nil
}
