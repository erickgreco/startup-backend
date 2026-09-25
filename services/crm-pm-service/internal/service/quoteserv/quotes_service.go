package quoteserv

import (
	"context"
	"fmt"
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/shared/idgen"
	"github.com/erickgreco/startup-backend/shared/logger"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	quotePrefix = "Q"
)

type QuoteService struct {
	quoteRepo domain.QuoteRepository
	leadRepo  domain.LeadRepository
}

func NewQuoteService(quoteRepo domain.QuoteRepository, leadRepo domain.LeadRepository) *QuoteService {
	return &QuoteService{
		quoteRepo: quoteRepo,
		leadRepo:  leadRepo,
	}
}

func (s *QuoteService) CreateQuote(ctx context.Context, leadID uuid.UUID, req *domain.QuoteRequest) (*domain.QuotesModel, error) {
	lead, err := s.leadRepo.GetLeadByID(ctx, leadID)
	if err != nil {
		return nil, fmt.Errorf("quote_service: create_quote: %w", err)
	}

	// TODO: logger solo manda en terminal, se necesita poder mostrar al usuario
	if lead.DataQuality == domain.Low {
		logger.Warn("warning", "data_quality: %v", lead.DataQuality)
	}

	if req.ProjectEstimatedDurationDays < 1 || req.ProjectEstimatedDurationDays > 180 {
		return nil, fmt.Errorf("quote_service: create_quote: %w", domain.ErrInvalidDuration)
	}

	if len(req.Items) == 0 {
		return nil, fmt.Errorf("quote_service: create_quote: %w", domain.ErrMissingItems)
	}

	quoteID := idgen.NewID()

	quoteItems := make([]domain.QuoteItem, len(req.Items))
	totalSubtotal := domain.ZeroMoney

	for i, item := range req.Items {
		if item.Quantity.LessThanOrEqual(decimal.Zero) || item.UnitPrice.LessThanOrEqual(domain.ZeroMoney) {
			return nil, fmt.Errorf("quote_service: create_quote: invalid_item: %d: %w", i, domain.ErrInvalidQuoteItem)
		}

		if item.UnitPrice.HasSubCentPrecision() {
			return nil, fmt.Errorf("quote_service: create_quote: invalid_item: %d:%w", i, domain.ErrInvalidUnitPrice)
		}

		subtotal := item.UnitPrice.Mul(item.Quantity)

		quoteItems[i] = domain.QuoteItem{
			ID:          idgen.NewID(),
			QuoteID:     quoteID,
			Description: item.Description,
			Quantity:    item.Quantity,
			Unit:        item.Unit,
			UnitPrice:   item.UnitPrice,
			SubTotal:    subtotal,
			SortOrder:   i,
			CreatedAt:   time.Now(),
		}

		totalSubtotal = totalSubtotal.Add(subtotal)
	}

	quoteNumber := idgen.GenerateCode(quotePrefix, quoteID)

	quote := &domain.QuotesModel{
		ID:                           quoteID,
		LeadID:                       leadID,
		QuoteNumber:                  quoteNumber,
		ClientName:                   req.ClientName,
		ClientBusiness:               req.ClientBusiness,
		ClientEmail:                  req.ClientEmail,
		ClientPhone:                  req.ClientPhone,
		ProjectTitle:                 req.ProjectTitle,
		ProjectDescription:           req.ProjectDescription,
		ProjectEstimatedDurationDays: req.ProjectEstimatedDurationDays,
		Items:                        quoteItems,
		PaymentMethod:                req.PaymentMethod,
		Status:                       domain.QuoteStatusDraft,
		CreatedAt:                    time.Now(),
		Subtotal:                     totalSubtotal, // ! Por ahora los campos son iguales, falta calcular taxes
		Total:                        totalSubtotal, // ! Por ahora los campos son iguales, falta calcular taxes
		Currency:                     req.Currency,
		Notes:                        req.Notes,
	}

	createdQuote, err := s.quoteRepo.CreateQuote(ctx, quote)
	if err != nil {
		return nil, fmt.Errorf("quote_service: create_quote: %w", err)
	}

	return createdQuote, nil
}

func (s *QuoteService) GetQuoteByID(ctx context.Context, id uuid.UUID) (*domain.QuotesModel, error) {
	quote, err := s.quoteRepo.GetQuoteByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("quote_service: get_quote_by_id: %w", err)
	}
	return quote, nil
}

// * Metodo de pago se valida en dto por conversion a dominio (string <-> domain)
// * Actualiza el draft, no cambia status ni fechas (solo updatedAt a nivel db)
func (s *QuoteService) UpdateQuote(ctx context.Context, quoteID uuid.UUID, req *domain.QuoteRequest) (*domain.QuotesModel, error) {
	if err := validateRequest(req); err != nil {
		return nil, fmt.Errorf("quote_service: create_manual_quote: %w", err)
	}

	existing, err := s.quoteRepo.GetQuoteByID(ctx, quoteID)
	if err != nil {
		return nil, fmt.Errorf("quote_service: update_quote: get_existing_quote: %w", err)
	}

	if existing.Status != domain.QuoteStatusDraft {
		return nil, domain.ErrInvalidStatusForUpdate
	}

	items := make([]domain.QuoteItem, len(req.Items))
	totalSubtotal := domain.ZeroMoney
	itemCreation := time.Now()

	for i, item := range req.Items {
		subtotal := item.UnitPrice.Mul(item.Quantity)

		items[i] = domain.QuoteItem{
			ID:          idgen.NewID(),
			QuoteID:     existing.ID,
			Description: item.Description,
			Quantity:    item.Quantity,
			Unit:        item.Unit,
			UnitPrice:   item.UnitPrice,
			SubTotal:    subtotal,
			SortOrder:   i,
			CreatedAt:   itemCreation,
		}

		totalSubtotal = totalSubtotal.Add(subtotal)
	}

	newQuote := &domain.QuotesModel{
		ID:                           existing.ID,
		LeadID:                       existing.LeadID,
		QuoteNumber:                  existing.QuoteNumber,
		ClientName:                   req.ClientName,
		ClientBusiness:               req.ClientBusiness,
		ClientEmail:                  req.ClientEmail,
		ClientPhone:                  req.ClientPhone,
		ProjectTitle:                 req.ProjectTitle,
		ProjectDescription:           req.ProjectDescription,
		ProjectResume:                req.ProjectResume,
		ProjectEstimatedDurationDays: req.ProjectEstimatedDurationDays,
		Items:                        items,
		FunctionalityExample:         req.FunctionalityExample,
		PaymentMethod:                req.PaymentMethod,
		Status:                       domain.QuoteStatusDraft,
		CreatedAt:                    existing.CreatedAt,
		Subtotal:                     totalSubtotal, // ! Por ahora los campos son iguales, falta calcular taxes
		Total:                        totalSubtotal, // ! Por ahora los campos son iguales, falta calcular taxes
		Currency:                     req.Currency,
		Notes:                        req.Notes,
	}

	updated, err := s.quoteRepo.UpdateQuote(ctx, newQuote)
	if err != nil {
		return nil, fmt.Errorf("quote_service: update_quote: %w", err)
	}

	return updated, nil
}
