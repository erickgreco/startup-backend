package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/repository/sqlgen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	queryTimeDuration = time.Second * 5
)

type LeadRepository struct {
	queries *sqlgen.Queries
}

func NewLeadRepository(pool *pgxpool.Pool) *LeadRepository {
	return &LeadRepository{
		queries: sqlgen.New(pool),
	}
}

func (r *LeadRepository) CreateLead(ctx context.Context, lead *domain.LeadsModel) (*domain.LeadsModel, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	row, err := r.queries.CreateLead(ctx, leadDomainToGen(lead))
	if err != nil {
		return nil, fmt.Errorf("leads_repo: create_lead: %w", translateWriteError(err))
	}

	return genToLeadDomain(row), nil
}

func (r *LeadRepository) GetLeadsByStatus(ctx context.Context, status domain.LeadStatus) ([]*domain.LeadsModel, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	rows, err := r.queries.GetLeadsByStatus(ctx, string(status))
	if err != nil {
		return nil, fmt.Errorf("leads_repo: get_leads_by_status: %w", translateReadError(err))
	}

	leads := make([]*domain.LeadsModel, 0, len(rows))
	for _, row := range rows {
		leads = append(leads, genToLeadDomain(row))
	}

	return leads, nil
}

func (r *LeadRepository) GetLeadByID(ctx context.Context, id uuid.UUID) (*domain.LeadsModel, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	row, err := r.queries.GetLeadByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("leads_repo: get_lead_by_id: %w", translateReadError(err))
	}

	return genToLeadDomain(row), nil
}

func (r *LeadRepository) UpdateLeadStatus(ctx context.Context, id uuid.UUID, status domain.LeadStatus) (*domain.LeadsModel, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	row, err := r.queries.UpdateLeadStatus(ctx, string(status), id)
	if err != nil {
		return nil, fmt.Errorf("leads_repo: update_lead_status: %w", translateUpdateError(err))
	}

	return genToLeadDomain(row), nil
}

func (r *LeadRepository) DeleteLead(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	affected, err := r.queries.DeleteLead(ctx, id)
	if err != nil {
		return fmt.Errorf("leads_repo: delete_lead: %w", translateWriteError(err))
	}

	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
