package dto

import (
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/google/uuid"
)

type QuoteResponse struct {
	ID                           uuid.UUID      `json:"id"`
	LeadID                       *uuid.UUID     `json:"lead_id"`
	QuoteNumber                  string         `json:"quote_number"`
	ClientName                   string         `json:"client_name"`
	ClientBusiness               *string        `json:"client_business"`
	ClientEmail                  string         `json:"client_email"`
	ClientPhone                  string         `json:"client_phone_number"`
	ProjectTitle                 string         `json:"project_title"`
	ProjectDescription           string         `json:"project_description"`
	ProjectResume                *string        `json:"project_resume"`
	ProjectEstimatedDurationDays int32          `json:"project_estimated_duration_days"`
	FunctionalityExample         *string        `json:"functionality_example"`
	Items                        []ItemResponse `json:"items"`
	PaymentMethod                string         `json:"payment_method"`
	Status                       string         `json:"status"`
	CreatedAt                    time.Time      `json:"created_at"`
	ValidUntil                   *time.Time     `json:"valid_until"`
	SentAt                       *time.Time     `json:"sent_at"`
	AcceptedAt                   *time.Time     `json:"accepted_at"`
	RejectedAt                   *time.Time     `json:"rejected_at"`
	Subtotal                     string         `json:"subtotal"`
	Total                        string         `json:"total"`
	Currency                     string         `json:"currency"`
	Notes                        *string        `json:"notes"`
	UpdatedAt                    time.Time      `json:"updated_at"`
}

type ItemResponse struct {
	ID          uuid.UUID `json:"item_id"`
	QuoteID     uuid.UUID `json:"quote_id"`
	Description string    `json:"item_description"`
	Quantity    string    `json:"quantity"`
	Unit        string    `json:"unit"`
	UnitPrice   string    `json:"unit_price"`
	Subtotal    string    `json:"subtotal"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func FromItemDomain(i *domain.QuoteItem) ItemResponse {
	return ItemResponse{
		ID:          i.ID,
		QuoteID:     i.QuoteID,
		Description: i.Description,
		Quantity:    i.Quantity.String(),
		Unit:        i.Unit,
		UnitPrice:   i.UnitPrice.String(),
		Subtotal:    i.SubTotal.String(),
		SortOrder:   i.SortOrder,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

func FromQuoteDomain(q *domain.QuotesModel) QuoteResponse {
	items := make([]ItemResponse, len(q.Items))

	for i, item := range q.Items {
		items[i] = FromItemDomain(&item)
	}

	return QuoteResponse{
		ID:                           q.ID,
		LeadID:                       q.LeadID,
		QuoteNumber:                  q.QuoteNumber,
		ClientName:                   q.ClientName,
		ClientBusiness:               q.ClientBusiness,
		ClientEmail:                  q.ClientEmail,
		ClientPhone:                  q.ClientPhone,
		ProjectTitle:                 q.ProjectTitle,
		ProjectDescription:           q.ProjectDescription,
		ProjectResume:                q.ProjectResume,
		ProjectEstimatedDurationDays: q.ProjectEstimatedDurationDays,
		FunctionalityExample:         q.FunctionalityExample,
		Items:                        items,
		PaymentMethod:                string(q.PaymentMethod),
		Status:                       string(q.Status),
		CreatedAt:                    q.CreatedAt,
		ValidUntil:                   q.ValidUntil,
		SentAt:                       q.SentAt,
		AcceptedAt:                   q.AcceptedAt,
		RejectedAt:                   q.RejectedAt,
		Subtotal:                     q.Subtotal.String(),
		Total:                        q.Total.String(),
		Currency:                     q.Currency.String(),
		Notes:                        q.Notes,
		UpdatedAt:                    q.UpdatedAt,
	}
}
