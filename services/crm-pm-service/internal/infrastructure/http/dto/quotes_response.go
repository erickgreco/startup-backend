package dto

import (
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/go-playground/locales/currency"
	"github.com/google/uuid"
)

type QuoteResponse struct {
	ID                           uuid.UUID          `json:"id"`
	LeadID                       uuid.UUID          `json:"lead_id"`
	QuoteNumber                  string             `json:"quote_number"`
	ClientName                   string             `json:"client_name"`
	ClientBusiness               *string            `json:"client_business"`
	ClientEmail                  string             `json:"client_email"`
	ClientPhone                  string             `json:"client_phone_number"`
	ProjectTitle                 string             `json:"project_title"`
	ProjectDescription           string             `json:"project_description"`
	ProjectResume                *string            `json:"project_resume"`
	ProjectEstimatedDurationDays int                `json:"project_estimated_duration_days"`
	FunctionalityExample         *string            `json:"functionality_example"`
	Items                        []domain.QuoteItem `json:"items"`
	PaymentMethod                domain.PayMethod   `json:"payment_method"`
	Status                       domain.QuoteStatus `json:"status"`
	CreatedAt                    time.Time          `json:"created_at"`
	ValidUntil                   *time.Time         `json:"valid_until"`
	SentAt                       *time.Time         `json:"sent_at"`
	AcceptedAt                   *time.Time         `json:"accepted_at"`
	RejectedAt                   *time.Time         `json:"rejected_at"`
	Subtotal                     domain.Money       `json:"subtotal"`
	Total                        domain.Money       `json:"total"`
	Currency                     currency.Type      `json:"currency"`
	Notes                        *string            `json:"notes"`
	UpdatedAt                    time.Time          `json:"updated_at"`
}

func FromQuoteDomain(q *domain.QuotesModel) QuoteResponse {
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
		Items:                        q.Items,
		PaymentMethod:                q.PaymentMethod,
		Status:                       q.Status,
		CreatedAt:                    q.CreatedAt,
		ValidUntil:                   q.ValidUntil,
		SentAt:                       q.SentAt,
		AcceptedAt:                   q.AcceptedAt,
		RejectedAt:                   q.RejectedAt,
		Subtotal:                     q.Subtotal,
		Total:                        q.Total,
		Currency:                     q.Currency,
		Notes:                        q.Notes,
		UpdatedAt:                    q.UpdatedAt,
	}
}
