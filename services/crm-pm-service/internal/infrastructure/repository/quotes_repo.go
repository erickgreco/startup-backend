package repository

import (
	"context"
	"fmt"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
// ** validUnit, sentAt, acceptedAt, rejectedAt
func (r *QuoteRepository) CreateQuote(ctx context.Context, quote *domain.QuotesModel) (*domain.QuotesModel, error) {
	var createdQuote domain.QuotesModel

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: create_quote: begin_transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	quoteQuery := `
		INSERT INTO quotes (
			id, lead_id, quote_number, client_name, client_business, client_email, client_phone,
			project_title, project_description, project_estimated_duration_days,
			payment_method, status, created_at, subtotal, total, 
			currency, notes
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17
		)
		RETURNING id, lead_id, quote_number, client_name, client_business, client_email, client_phone,
			project_title, project_description, project_estimated_duration_days,
			payment_method, status, created_at, subtotal, total, 
			currency, notes, updated_at
	`

	err = tx.QueryRow(
		ctx,
		quoteQuery,
		quote.ID,
		quote.LeadID,
		quote.QuoteNumber,
		quote.ClientName,
		quote.ClientBusiness,
		quote.ClientEmail,
		quote.ClientPhone,
		quote.ProjectTitle,
		quote.ProjectDescription,
		quote.ProjectEstimatedDurationDays,
		quote.PaymentMethod,
		quote.Status,
		quote.CreatedAt,
		quote.Subtotal,
		quote.Total,
		quote.Currency,
		quote.Notes,
	).Scan(
		&createdQuote.ID,
		&createdQuote.LeadID,
		&createdQuote.QuoteNumber,
		&createdQuote.ClientName,
		&createdQuote.ClientBusiness,
		&createdQuote.ClientEmail,
		&createdQuote.ClientPhone,
		&createdQuote.ProjectTitle,
		&createdQuote.ProjectDescription,
		&createdQuote.ProjectEstimatedDurationDays,
		&createdQuote.PaymentMethod,
		&createdQuote.Status,
		&createdQuote.CreatedAt,
		&createdQuote.Subtotal,
		&createdQuote.Total,
		&createdQuote.Currency,
		&createdQuote.Notes,
		&createdQuote.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: create_quote: insert_quote: %w", err)
	}

	if len(quote.Items) > 0 {
		itemsBatch := &pgx.Batch{}
		itemsQuery := `
			INSERT INTO quote_items (
				id, quote_id, description, quantity, unit, unit_price, subtotal, sort_order,
				created_at
			)
			VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9
			)
			RETURNING id, quote_id, description, quantity, unit, unit_price, subtotal, sort_order,
				created_at, updated_at
		`

		for i, item := range quote.Items {
			itemsBatch.Queue(
				itemsQuery,
				item.ID,
				quote.ID,
				item.Description,
				item.Quantity,
				item.Unit,
				item.UnitPrice,
				item.SubTotal,
				i,
				item.CreatedAt,
			)
		}

		br := tx.SendBatch(ctx, itemsBatch)
		defer br.Close()

		for range quote.Items {
			var item domain.QuoteItem

			err := br.QueryRow().Scan(
				&item.ID,
				&item.QuoteID,
				&item.Description,
				&item.Quantity,
				&item.Unit,
				&item.UnitPrice,
				&item.SubTotal,
				&item.SortOrder,
				&item.CreatedAt,
				&item.UpdatedAt,
			)
			if err != nil {
				return nil, fmt.Errorf("quotes_repo: create_quote: insert_quote_items: %w", err)
			}

			createdQuote.Items = append(createdQuote.Items, item)
		}

		if err := br.Close(); err != nil {
			return nil, fmt.Errorf("quotes_repo: create_quote: close_batch: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("quotes_repo: create_quote: commit_transaction: %w", err)
	}

	return &createdQuote, nil
}

func (r *QuoteRepository) GetQuoteByID(ctx context.Context, id uuid.UUID) (*domain.QuotesModel, error) {
	var quote domain.QuotesModel

	quoteQuery := `
		SELECT id, lead_id, quote_number, client_name, client_business, client_email, client_phone,
			project_title, project_description, project_resume, project_estimated_duration_days,
			functionality_example, payment_method, status, created_at, valid_until, sent_at,
			accepted_at, rejected_at, subtotal, total, currency, notes, updated_at
		FROM quotes
		WHERE id = $1
	`

	err := r.pool.QueryRow(
		ctx,
		quoteQuery,
		id,
	).Scan(
		&quote.ID,
		&quote.LeadID,
		&quote.QuoteNumber,
		&quote.ClientName,
		&quote.ClientBusiness,
		&quote.ClientEmail,
		&quote.ClientPhone,
		&quote.ProjectTitle,
		&quote.ProjectDescription,
		&quote.ProjectResume,
		&quote.ProjectEstimatedDurationDays,
		&quote.FunctionalityExample,
		&quote.PaymentMethod,
		&quote.Status,
		&quote.CreatedAt,
		&quote.ValidUntil,
		&quote.SentAt,
		&quote.AcceptedAt,
		&quote.RejectedAt,
		&quote.Subtotal,
		&quote.Total,
		&quote.Currency,
		&quote.Notes,
		&quote.UpdatedAt,
	)
	if err != nil {
		mappedErr := mapReadError(err)
		return nil, fmt.Errorf("quotes_repo: get_quote_by_id: %w", mappedErr)
	}

	itemsQuery := `
		SELECT id, quote_id, description, quantity, unit_price, unit, subtotal, sort_order,
			created_at, updated_at
		FROM quote_items
		WHERE quote_id = $1
		ORDER BY sort_order ASC
	`

	rows, err := r.pool.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: get_quote_by_id: select_quote_items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.QuoteItem

		err := rows.Scan(
			&item.ID,
			&item.QuoteID,
			&item.Description,
			&item.Quantity,
			&item.UnitPrice,
			&item.Unit,
			&item.SubTotal,
			&item.SortOrder,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("quotes_repo: get_quote_by_id: scan_quote_item: %w", err)
		}

		quote.Items = append(quote.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("quotes_repo: get_quote_by_id: iterate_quote_items: %w", err)
	}

	return &quote, nil
}
