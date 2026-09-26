package repository

import (
	"context"
	"fmt"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/repository/sqlgen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuoteRepository struct {
	pool    *pgxpool.Pool
	queries *sqlgen.Queries
}

func NewQuotesRepository(pool *pgxpool.Pool) *QuoteRepository {
	return &QuoteRepository{
		pool:    pool,
		queries: sqlgen.New(pool),
	}
}

// ** Create asume que solo se esta creando un draft por lo que deja fuera los campos
// ** validUnit, sentAt, acceptedAt, rejectedAt
func (r *QuoteRepository) CreateQuote(ctx context.Context, quote *domain.QuotesModel) (*domain.QuotesModel, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: create_quote: begin_transaction: %w", translateWriteError(err))
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	quoteRow, err := qtx.CreateQuote(ctx, quoteDomainToGen(quote))
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: create_quote: insert_quote: %w", translateWriteError(err))
	}

	itemRows := make([]sqlgen.QuoteItem, 0, len(quote.Items))
	if len(quote.Items) > 0 {
		var batchErr error

		results := qtx.CreateQuoteItems(ctx, quoteItemsDomainToGen(quote.Items))
		results.QueryRow(func(_ int, qi sqlgen.QuoteItem, err error) {
			if err != nil {
				if batchErr == nil {
					batchErr = err
				}
				return
			}
			itemRows = append(itemRows, qi)
		})
		if err := results.Close(); err != nil && batchErr == nil {
			batchErr = err
		}

		if batchErr != nil {
			return nil, fmt.Errorf("quotes_repo: create_quote: insert_quote_items: %w", translateWriteError(batchErr))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("quotes_repo: create_quote: commit_transaction: %w", translateWriteError(err))
	}

	created, err := genToQuoteDomain(quoteRow, itemRows)
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: create_quote: map_result: %w", fmt.Errorf("%w: %w", domain.ErrInternal, err))
	}

	return created, nil
}

func (r *QuoteRepository) GetQuoteByID(ctx context.Context, id uuid.UUID) (*domain.QuotesModel, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	quoteRow, err := r.queries.GetQuoteByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: get_quote_by_id: %w", translateReadError(err))
	}

	itemRows, err := r.queries.ListQuoteItemsByQuoteID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: get_quote_by_id: list_items: %w", translateReadError(err))
	}

	quote, err := genToQuoteDomain(quoteRow, itemRows)
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: get_quote_by_id: map_result: %w", translateReadError(err))
	}

	return quote, nil
}

func (r *QuoteRepository) UpdateQuote(ctx context.Context, quote *domain.QuotesModel) (*domain.QuotesModel, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: update_quote: begin_transaction: %w", translateWriteError(err))
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	quoteRow, err := qtx.UpdateQuote(ctx, quoteDomainToGenUpdate(quote))
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: update_quote: update_quote: %w", translateUpdateError(err))
	}

	affected, err := qtx.DeleteQuoteItemsByQuoteID(ctx, quote.ID)
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: update_quote: delete_quote_items: %w", translateWriteError(err))
	}

	if affected == 0 {
		return nil, domain.ErrInternal
	}

	itemRows := make([]sqlgen.QuoteItem, 0, len(quote.Items))
	if len(quote.Items) > 0 {
		var batchErr error

		results := qtx.CreateQuoteItems(ctx, quoteItemsDomainToGen(quote.Items))
		results.QueryRow(func(_ int, qi sqlgen.QuoteItem, err error) {
			if err != nil {
				if batchErr == nil {
					batchErr = err
				}
				return
			}
			itemRows = append(itemRows, qi)
		})
		if err := results.Close(); err != nil && batchErr == nil {
			batchErr = err
		}

		if batchErr != nil {
			return nil, fmt.Errorf("quotes_repo: update_quote: insert_quote_items: %w", translateWriteError(batchErr))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("quotes_repo: update_quote: commit_transaction: %w", translateWriteError(err))
	}

	updated, err := genToQuoteDomain(quoteRow, itemRows)
	if err != nil {
		return nil, fmt.Errorf("quotes_repo: update_quote: map_result: %w", fmt.Errorf("%w: %w", domain.ErrInternal, err))
	}

	return updated, nil
}
