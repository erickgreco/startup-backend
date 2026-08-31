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
	`

	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	_, err := r.pool.Exec(
		ctx,
		query,
		lead.ID,
		lead.Name,
		lead.Email,
		lead.Message,
		lead.Status,
		lead.Source,
	)
	if err != nil {
		return err
	}

	return nil
}
