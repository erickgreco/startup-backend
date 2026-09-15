package quoteserv

import (
	"context"
	"fmt"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/shared/logger"
	"github.com/google/uuid"
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

func (s QuoteService) CreateQuote(ctx context.Context, id uuid.UUID, req *domain.QuoteRequest) (*domain.QuotesModel, error) {
	// * Obtener lead
	lead, err := s.leadRepo.GetLeadByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("quote_service: create_quote: %w", err)
	}

	// * Validar calidad de los datos
	if lead.DataQuality == domain.Low {
		logger.Warn("warning: data_quality: %v", lead.DataQuality)
	}

	// * Validar que contiene items (minimo 1)
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("quote_service: create_quote: %w", domain.ErrMissingItems)
	}

	// TODO: validar metodo de pago (viene del req)

	if !req.PaymentMethod.IsValid() {
		return nil, fmt.Errorf("quote_service: create_quote: %w", domain.ErrPaymethodRequired)
	}

	// TODO: deserializar items, validar subtotal == total

	// TODO: crear quote number (crear helper func)

	// TODO: crear draft previo a persistencia db

	// TODO: persistir a db

	return nil, nil
}
