package repository

import (
	"context"
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
		INSERT INTO Leads (
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
		return mapCreateLeadDBError(err)
	}

	return nil
}
