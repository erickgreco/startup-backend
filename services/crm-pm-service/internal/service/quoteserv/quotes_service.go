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

// TODO: currency se esta almacenando como int
func (s *QuoteService) CreateQuote(ctx context.Context, leadID uuid.UUID, req *domain.QuoteRequest) (*domain.QuotesModel, error) {
	lead, err := s.leadRepo.GetLeadByID(ctx, leadID)
	if err != nil {
		return nil, fmt.Errorf("quote_service: create_quote: %w", err)
	}

	if lead.DataQuality == domain.Low {
		logger.Warn("warning: data_quality: %v", lead.DataQuality)
	}

	if req.ProjectEstimatedDurationDays < 1 || req.ProjectEstimatedDurationDays > 180 {
		return nil, fmt.Errorf("quote_service: create_quote: %w", domain.ErrInvalidDuration)
	}

	if len(req.Items) == 0 {
		return nil, fmt.Errorf("quote_service: create_quote: %w", domain.ErrMissingItems)
	}

	if !req.PaymentMethod.IsValid() {
		return nil, fmt.Errorf("quote_service: create_quote: %w", domain.ErrPaymethodRequired)
	}

	quoteID := idgen.NewID()

	quoteItems := make([]domain.QuoteItem, len(req.Items))
	toalSubtotal := decimal.Zero

	for i, item := range req.Items {
		if item.Quantity.LessThanOrEqual(decimal.Zero) || item.UnitPrice.LessThanOrEqual(decimal.Zero) {
			return nil, fmt.Errorf("quote_service: create_quote: invalid_item: %d: %w", i, domain.ErrInvalidQuoteItem)
		}

		subtotal := item.Quantity.Mul(item.UnitPrice)

		quoteItems[i] = domain.QuoteItem{
			ID:          idgen.NewID(),
			QuoteID:     quoteID,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			SubTotal:    subtotal,
			SortOrder:   i,
		}

		toalSubtotal = toalSubtotal.Add(subtotal)
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
		Subtotal:                     toalSubtotal, // ! Por ahora los campos son iguales, falta calcular taxes
		Total:                        toalSubtotal, // ! Por ahora los campos son iguales, falta calcular taxes
		Currency:                     req.Currency,
		Notes:                        req.Notes,
	}

	createdQuote, err := s.quoteRepo.CreateQuote(ctx, quote)
	if err != nil {
		return nil, fmt.Errorf("quote_service: create_quote: %w", err)
	}

	return createdQuote, nil
}
