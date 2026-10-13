package quoteserv

import (
	"os"
	"testing"
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/shopspring/decimal"
)

func TestRenderQuotePDF(t *testing.T) {
	now := time.Now()
	validUntil := now.AddDate(0, 0, 30)
	website := "https://erickgreco.dev"
	business := "Greco Dev Studio"
	desc := "Desarrollo de plataforma web completa con dashboard administrativo"
	resume := "Sistema construido en Go con PostgreSQL, API REST y panel en React"
	duration := int32(90)
	funcExample := "Ejemplo de funcionalidad va aqui"
	cur, err := domain.ParseCurrency("MXN")
	if err != nil {
		t.Fatal(err)
	}

	up := &domain.UserProfileModel{
		FirstName:    "Erick",
		LastName:     "Gutierrez",
		Email:        "erick@example.com",
		PhoneNumber:  "+52 55 1234 5678",
		BusinessName: &business,
		WebsiteURL:   &website,
		Language:     "en",
	}

	q := &domain.QuotesModel{
		QuoteNumber:                  "QT-2026-001",
		ClientName:                   "Carlos Lopez",
		ClientEmail:                  "carlos@example.com",
		ClientPhone:                  "+52 55 9876 5432",
		ProjectTitle:                 "Plataforma SaaS para gestión de inventarios",
		ProjectDescription:           desc,
		ProjectResume:                &resume,
		ProjectEstimatedDurationDays: duration,
		FunctionalityExample:         &funcExample,
		Items: []domain.QuoteItem{
			{
				Description: "Diseno UI/UX del dashboard",
				Quantity:    decimal.NewFromInt(1),
				Unit:        "Proyecto",
				UnitPrice:   domain.FromCents(1500000),
				SubTotal:    domain.FromCents(1500000),
			},
			{
				Description: "Desarrollo backend API REST",
				Quantity:    decimal.NewFromInt(120),
				Unit:        "Horas",
				UnitPrice:   domain.FromCents(15000),
				SubTotal:    domain.FromCents(1800000),
			},
			{
				Description: "Integracion con pasarela de pagos",
				Quantity:    decimal.NewFromInt(1),
				Unit:        "Modulo",
				UnitPrice:   domain.FromCents(1200000),
				SubTotal:    domain.FromCents(1200000),
			},
		},
		PaymentMethod: domain.PaymentMilestone,
		Currency:      cur,
		CreatedAt:     now,
		ValidUntil:    &validUntil,
		Subtotal:      domain.FromCents(4500000),
		Total:         domain.FromCents(4500000),
	}

	pdf, err := RenderQuotePDF(q, up)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("test_quote.pdf", pdf, 0644)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("PDF generado: test_quote.pdf")
}
