package repository

import (
	"context"
	"fmt"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuoteRepository struct {
	pool *pgxpool.Pool
}

func NewQuotesRepository(pool *pgxpool.Pool) *QuoteRepository {
	return &QuoteRepository{
		pool: pool,
	}
}

// ** Create asume que solo se esta creando un draft por lo que deja fuera los campos
// ** sentAt, acceptedAt, rejectedAt
func (r *QuoteRepository) CreateQuote(ctx context.Context, quote *domain.QuotesModel) (*domain.QuotesModel, error) {
	var createdQuote domain.QuotesModel

	query := `
		INSERT INTO quotes (
			id, lead_id, quote_number, client_name, client_business, client_email, client_phone,
			project_title, project_description, project_resume, project_estimated_duration,
			functionality_example, quote_items, payment_method, status, created_at, valid_until,
			subtotal_cents, total_cents, currency, notes
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21
		)
		RETURNING id, lead_id, quote_number, client_name, client_business, client_email, client_phone,
			project_title, project_description, project_resume, project_estimated_duration,
			functionality_example, quote_items, payment_method, status, created_at, valid_until,
			subtotal_cents, total_cents, currency, notes, updated_at
	`

	itemstoJSON, err := quote.SetItemsToJSON()
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: create_quote: failed_to_serialize_items: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	row := r.pool.QueryRow(
		ctx,
		query,
		quote.ID,
		quote.LeadID,
		quote.QuoteNumber,
		quote.ClientName,
		quote.ClientBusiness,
		quote.ClientEmail,
		quote.ClientPhone,
		quote.ProjectTitle,
		quote.ProjectDescription,
		quote.ProjectResume,
		quote.ProjectEstimatedDurationSec,
		quote.FunctionalityExample,
		itemstoJSON,
		quote.PaymentMethod,
		quote.Status,
		quote.CreatedAt,
		quote.ValidUntil,
		quote.Subtotal,
		quote.Total,
		quote.Currency,
		quote.Notes,
	)

	var JSONToItems []byte

	err = row.Scan(
		&createdQuote.ID,
		&createdQuote.LeadID,
		&createdQuote.QuoteNumber,
		&createdQuote.ClientName,
		&createdQuote.ClientBusiness,
		&createdQuote.ClientEmail,
		&createdQuote.ClientPhone,
		&createdQuote.ProjectTitle,
		&createdQuote.ProjectDescription,
		&createdQuote.ProjectResume,
		&createdQuote.ProjectEstimatedDurationSec,
		&createdQuote.FunctionalityExample,
		&JSONToItems,
		&createdQuote.PaymentMethod,
		&createdQuote.Status,
		&createdQuote.CreatedAt,
		&createdQuote.ValidUntil,
		&createdQuote.Subtotal,
		&createdQuote.Total,
		&createdQuote.Currency,
		&createdQuote.Notes,
		&createdQuote.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: create_quote: %w", err)
	}

	if err := quote.SetItemsFromJSON(JSONToItems); err != nil {
		return nil, fmt.Errorf("quotes_repo: create_quote: failed_to_deserialize_items: %w", err)
	}

	return &createdQuote, nil
}
