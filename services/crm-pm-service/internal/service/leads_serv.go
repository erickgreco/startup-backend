package service

import (
	"context"
	"fmt"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/google/uuid"
)

type LeadService struct {
	repo domain.LeadRepository
}

func NewLeadService(repo domain.LeadRepository) *LeadService {
	return &LeadService{
		repo: repo,
	}
}

func (s *LeadService) CreateLead(ctx context.Context, lead *domain.LeadsModel) error {
	if !lead.Source.IsValid() {
		return fmt.Errorf("service: create_lead: source %q: %w", lead.Source, domain.ErrInvalidSource)
	}

	if lead.Name == "" || lead.Email == "" {
		return fmt.Errorf("service: create_lead: %w", domain.ErrMissingRequiredFields)
	}

	lead.ID = uuid.New()
	lead.Status = domain.LeadStatusNew

	if err := s.repo.CreateLead(ctx, lead); err != nil {
		return fmt.Errorf("service: created_lead: %w", err)
	}

	return nil
}
