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
		return fmt.Errorf("leads_service: create_lead: source %q: %w", lead.Source, domain.ErrInvalidSource)
	}

	if lead.Name == "" || lead.Email == "" {
		return fmt.Errorf("leads_service: create_lead: %w", domain.ErrMissingRequiredFields)
	}

	lead.ID = uuid.New()
	lead.Status = domain.LeadStatusNew

	if err := s.repo.CreateLead(ctx, lead); err != nil {
		return fmt.Errorf("leads_service: created_lead: %w", err)
	}

	return nil
}

func (s *LeadService) GetLeadsByStatus(ctx context.Context, status domain.LeadStatus) ([]*domain.LeadsModel, error) {
	leads, err := s.repo.GetLeadsByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("leads_service: get_lead_by_status: %w", err)
	}

	return leads, nil
}

func (s *LeadService) GetLeadByID(ctx context.Context, id uuid.UUID) (*domain.LeadsModel, error) {
	lead, err := s.repo.GetLeadByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("leads_service: get_lead_by_id: %w", err)
	}

	return lead, nil
}

func (s *LeadService) UpdateLeadStatus(ctx context.Context, id uuid.UUID, status domain.LeadStatus) (*domain.LeadsModel, error) {
	if !status.IsValid() {
		return nil, fmt.Errorf("leads_service: update_lead_status: status: %q: %w", status, domain.ErrInvalidData)
	}

	lead, err := s.repo.UpdateLeadStatus(ctx, id, status)
	if err != nil {
		return nil, fmt.Errorf("leads_service: update_lead_status: %w", err)
	}

	return lead, nil
}

func (s *LeadService) DeleteLead(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteLead(ctx, id); err != nil {
		return fmt.Errorf("leads_service: delete_lead: %w", err)
	}

	return nil
}
