package leadserv

import (
	"context"
	"fmt"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/shared/idgen"
	"github.com/erickgreco/startup-backend/shared/logger"
	"github.com/google/uuid"
)

type LeadService struct {
	repo     domain.LeadRepository
	analyzer *LocalLeadAnalyzer
}

func NewLeadService(repo domain.LeadRepository) *LeadService {
	return &LeadService{
		repo:     repo,
		analyzer: NewLocalLeadAnalyzer(),
	}
}

func (s *LeadService) CreateLead(ctx context.Context, req *domain.LeadRequest) (*domain.LeadsModel, error) {
	if !req.Source.IsValid() {
		return nil, fmt.Errorf("leads_service: create_lead: source %q: %w", req.Source, domain.ErrInvalidSource)
	}

	if req.Name == "" || req.Email == "" {
		return nil, fmt.Errorf("leads_service: create_lead: %w", domain.ErrMissingRequiredFields)
	}

	analysis := s.analyzer.Analyze(req.Message)

	logger.Info("lead_analysis", "score:", analysis.Score, "quality:", analysis.QA, "concepts:", analysis.FoundConcepts)

	model := &domain.LeadsModel{
		ID:           idgen.NewID(),
		Name:         req.Name,
		BusinessName: req.BusinessName,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		Message:      req.Message,
		DataQuality:  analysis.QA,
		QualityScore: analysis.Score,
		Status:       domain.LeadStatusNew,
		Source:       req.Source,
	}

	lead, err := s.repo.CreateLead(ctx, model)
	if err != nil {
		return nil, fmt.Errorf("leads_service: create_lead: %w", err)
	}

	return lead, nil
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
