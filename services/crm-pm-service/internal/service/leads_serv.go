package service

import (
	"context"
	"errors"
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

// TODO: trabajar con el source
func (s *LeadService) CreateLead(ctx context.Context, lead *domain.LeadsModel) error {
	if !lead.Source.IsValid() {
		return fmt.Errorf("leads_service: invalid_source %q: %w", lead.Source, errors.New("invalid_input")) //TODO: centralizar errores
	}

	if lead.Name == "" || lead.Email == "" {
		return errors.New("name and email are required")
	}

	lead.ID = uuid.New()
	lead.Status = "new"

	return s.repo.CreateLead(ctx, lead)
}
