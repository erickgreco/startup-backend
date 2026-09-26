package dto

import (
	"fmt"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/shopspring/decimal"
)

type CreateQuoteRequest struct {
	ClientName                   string        `json:"client_name" binding:"required"`
	ClientBusiness               *string       `json:"client_business,omitempty"`
	ClientEmail                  string        `json:"client_email" binding:"required"`
	ClientPhone                  string        `json:"client_phone_number" binding:"required"`
	ProjectTitle                 string        `json:"project_title" binding:"required"`
	ProjectDescription           string        `json:"project_description" binding:"required"`
	ProjectResume                *string       `json:"project_resume,omitempty"`
	ProjectEstimatedDurationDays int           `json:"project_estimated_duration_days" binding:"required"`
	Items                        []ItemRequest `json:"items" binding:"required"`
	FunctionalityExample         *string       `json:"functionality_example,omitempty"`
	PaymentMethod                string        `json:"payment_method" binding:"required"`
	Currency                     string        `json:"currency" binding:"required"`
	Notes                        *string       `json:"notes,omitempty"`
}

func (r CreateQuoteRequest) ToDomain() (*domain.QuoteRequest, error) {
	cur, err := domain.ParseCurrency(r.Currency)
	if err != nil {
		return nil, fmt.Errorf("dto: invalid_currency: %w", err)
	}

	payMethod := domain.PayMethod(r.PaymentMethod)
	if !payMethod.IsValid() {
		return nil, fmt.Errorf("dto: invalid_payment_method: %q", r.PaymentMethod)
	}

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
		ProjectResume:                r.ProjectResume,
		ProjectEstimatedDurationDays: r.ProjectEstimatedDurationDays,
		Items:                        items,
		FunctionalityExample:         r.FunctionalityExample,
		PaymentMethod:                payMethod,
		Currency:                     cur,
		Notes:                        r.Notes,
	}, nil
}

type ItemRequest struct {
	Description string          `json:"description" binding:"required"`
	Quantity    decimal.Decimal `json:"quantity" binding:"required"`
	Unit        string          `json:"unit" binding:"required"`
	UnitPrice   decimal.Decimal `json:"unit_price" binding:"required"`
}

func (i ItemRequest) itemToDomain() domain.ItemRequest {
	return domain.ItemRequest{
		Description: i.Description,
		Quantity:    i.Quantity,
		Unit:        i.Unit,
		UnitPrice:   domain.Money(i.UnitPrice),
	}
}
