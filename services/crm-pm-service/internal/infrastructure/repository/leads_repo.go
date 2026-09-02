package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/google/uuid"
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

func (r *LeadRepository) GetLeadsByStatus(ctx context.Context, status domain.LeadStatus) ([]*domain.LeadsModel, error) {
	query := `
		SELECT id, name, email, message, status, source, created_at, updated_at
		FROM leads
		WHERE status = $1
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	rows, err := r.pool.Query(
		ctx,
		query,
		status,
	)
	if err != nil {
		mappedErr := mapReadError(err)
		return nil, fmt.Errorf("leads_repo: get_lead_by_status: %w", mappedErr)
	}
	defer rows.Close()

	var leads []*domain.LeadsModel
	for rows.Next() {
		var lead domain.LeadsModel
		if err := rows.Scan(
			&lead.ID,
			&lead.Name,
			&lead.Email,
			&lead.Message,
			&lead.Status,
			&lead.Source,
			&lead.CreatedAt,
			&lead.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("leads_repo: get_leads_by_status: %w", err)
		}
		leads = append(leads, &lead)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("leads_repo: get_leads_by_status: %w", err)
	}

	if len(leads) == 0 {
		return nil, domain.ErrNotFound
	}

	return leads, nil
}

func (r *LeadRepository) GetLeadByID(ctx context.Context, id uuid.UUID) (*domain.LeadsModel, error) {
	var lead domain.LeadsModel

	query := `
		SELECT id, name, email, message, status, source, created_at, updated_at
		FROM leads
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	err := r.pool.QueryRow(
		ctx,
		query,
		id,
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
		return nil, fmt.Errorf("leads_repo: get_lead_by_id: %w", mappedErr)
	}

	return &lead, nil
}
