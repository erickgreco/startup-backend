package dto

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/go-playground/locales/currency"
	"github.com/shopspring/decimal"
)

type CreateQuoteRequest struct {
	ClientName                   string           `json:"client_name" binding:"required"`
	ClientBusiness               *string          `json:"client_business,omitempty"`
	ClientEmail                  string           `json:"client_email" binding:"required"`
	ClientPhone                  string           `json:"client_phone_number" binding:"required"`
	ProjectTitle                 string           `json:"project_title" binding:"required"`
	ProjectDescription           string           `json:"project_description" binding:"required"`
	ProjectEstimatedDurationDays int              `json:"project_estimated_duration_days" binding:"required"`
	Items                        []ItemRequest    `json:"items" binding:"required"`
	PaymentMethod                domain.PayMethod `json:"payment_method" binding:"required"`
	Currency                     currency.Type    `json:"currency" binding:"required"`
	Notes                        *string          `json:"notes,omitempty"`
}

func (r CreateQuoteRequest) ToDomain() *domain.QuoteRequest {
	items := make([]domain.ItemRequest, len(r.Items))
	for i, item := range r.Items {
		items[i] = item.itemToDomain()
	}

	return &domain.QuoteRequest{
		ClientName:                   r.ClientName,
		ClientBusiness:               r.ClientBusiness,
		ClientEmail:                  r.ClientEmail,
		ClientPhone:                  r.ClientPhone,
		ProjectTitle:                 r.ProjectTitle,
		ProjectDescription:           r.ProjectDescription,
		ProjectEstimatedDurationDays: r.ProjectEstimatedDurationDays,
		Items:                        items,
		PaymentMethod:                r.PaymentMethod,
		Currency:                     r.Currency,
		Notes:                        r.Notes,
	}
}

type ItemRequest struct {
	Description string          `json:"description" binding:"required"`
	Quantity    decimal.Decimal `json:"quantity" binding:"required"`
	UnitPrice   decimal.Decimal `json:"unit_price" binding:"required"`
}

func (i ItemRequest) itemToDomain() domain.ItemRequest {
	return domain.ItemRequest{
		Description: i.Description,
		Quantity:    i.Quantity,
		UnitPrice:   i.UnitPrice,
	}
}
