package dto

import (
	"fmt"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/shopspring/decimal"
)

type CreateQuoteRequest struct {
	ClientName                   string        `json:"client_name" binding:"required"`
	ClientBusiness               *string       `json:"client_business,omitempty"`
	ClientEmail                  string        `json:"client_email" binding:"required,email"`
	ClientPhone                  string        `json:"client_phone_number" binding:"required"`
	ProjectTitle                 string        `json:"project_title" binding:"required"`
	ProjectDescription           string        `json:"project_description" binding:"required"`
	ProjectResume                *string       `json:"project_resume,omitempty"`
	ProjectEstimatedDurationDays int32         `json:"project_estimated_duration_days" binding:"required,gt=0"`
	Items                        []ItemRequest `json:"items" binding:"required,dive"`
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
		iDom, err := item.itemToDomain()
		if err != nil {
			return nil, fmt.Errorf("dto: items[%d]: %w", i, err)
		}

		items[i] = *iDom
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
	Description string `json:"description" binding:"required"`
	Quantity    string `json:"quantity" binding:"required"`
	Unit        string `json:"unit" binding:"required"`
	UnitPrice   string `json:"unit_price" binding:"required"`
}

func (i ItemRequest) itemToDomain() (*domain.ItemRequest, error) {
	qty, err := decimal.NewFromString(i.Quantity)
	if err != nil {
		return nil, fmt.Errorf("dto: invalid_quantity: %q: %w", i.Quantity, err)
	}

	price, err := domain.ParseMoney(i.UnitPrice)
	if err != nil {
		return nil, fmt.Errorf("dto: invalid_unit_price: %w", err)
	}

	return &domain.ItemRequest{
		Description: i.Description,
		Quantity:    qty,
		Unit:        i.Unit,
		UnitPrice:   price,
	}, nil
}
