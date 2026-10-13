package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type QuoteRequest struct {
	ClientName                   string
	ClientBusiness               *string
	ClientEmail                  string
	ClientPhone                  string
	ProjectTitle                 string
	ProjectDescription           string
	ProjectResume                *string
	ProjectEstimatedDurationDays int32
	FunctionalityExample         *string
	Items                        []ItemRequest
	PaymentMethod                PayMethod
	Currency                     Currency // TODO: trabajar con location una vez se implemente middleware
	Notes                        *string
}

type QuotesModel struct {
	ID                           uuid.UUID
	LeadID                       *uuid.UUID
	QuoteNumber                  string
	ClientName                   string
	ClientBusiness               *string
	ClientEmail                  string
	ClientPhone                  string
	ProjectTitle                 string
	ProjectDescription           string
	ProjectResume                *string
	ProjectEstimatedDurationDays int32
	FunctionalityExample         *string
	Items                        []QuoteItem
	PaymentMethod                PayMethod
	Status                       QuoteStatus
	CreatedAt                    time.Time
	ValidUntil                   *time.Time
	SentAt                       *time.Time
	AcceptedAt                   *time.Time
	RejectedAt                   *time.Time
	Subtotal                     Money
	Total                        Money
	Currency                     Currency
	Notes                        *string
	UpdatedAt                    time.Time
}

type PayMethod string

const (
	/*
		- Proyectos pequenos
		- Clientes nuevos sin track record
		- Trabajos de bajo riesgo (plantillas o disenos simples)
	*/
	PaymentUpfront PayMethod = "Upfront"
	/*
		- Desarrollo complejo
		- Clientes corporativos
		- Relaciones a largo plazo
	*/
	PaymentMilestone PayMethod = "Milestone"
	/*
		- Proyectos pequenos
		- Clientes con historial establecido
		- Entregables simples/tangibles
		- NO RECOMENDADO PARA PROYECTOS GRANDES
	*/
	PaymentUponDelivery PayMethod = "Upon delivery"
	/*
		- Proyectos costosos
		- Clientes nuevos pero con referencias
		- Proteccion sin ahuyentar al cliente
	*/
	PaymentSplit PayMethod = "Split"
)

func (p PayMethod) IsValid() bool {
	switch p {
	case PaymentUpfront, PaymentMilestone, PaymentUponDelivery, PaymentSplit:
		return true
	default:
		return false
	}
}

func (p PayMethod) Translate(lang string) string {
	if lang == "es" {
		switch p {
		case PaymentUpfront:
			return "Pago anticipado"
		case PaymentMilestone:
			return "Por entregable"
		case PaymentUponDelivery:
			return "Contra entrega"
		case PaymentSplit:
			return "Pago dividido"
		default:
			return string(p)
		}
	}

	return string(p)
}

type QuoteStatus string

const (
	QuoteStatusDraft     QuoteStatus = "draft"
	QuoteStatusSent      QuoteStatus = "sent"
	QuoteStatusAccepted  QuoteStatus = "accepted"
	QuoteStatusRejected  QuoteStatus = "rejected"
	QuoteStatusExpired   QuoteStatus = "expired"
	QuoteStatusConverted QuoteStatus = "converted"
)

func (s QuoteStatus) IsValid() bool {
	switch s {
	case QuoteStatusDraft, QuoteStatusSent, QuoteStatusAccepted, QuoteStatusRejected, QuoteStatusExpired, QuoteStatusConverted:
		return true
	default:
		return false
	}
}

type QuoteRepository interface {
	CreateQuote(ctx context.Context, quote *QuotesModel) (*QuotesModel, error)
	GetQuoteByID(ctx context.Context, id uuid.UUID) (*QuotesModel, error)
	UpdateQuote(ctx context.Context, quote *QuotesModel) (*QuotesModel, error)
	DeleteQuote(ctx context.Context, id uuid.UUID) error
}

type QuoteService interface {
	CreateManualQuote(ctx context.Context, req *QuoteRequest) (*QuotesModel, error)
	CreateQuoteFromLead(ctx context.Context, leadID uuid.UUID, req *QuoteRequest) (*QuotesModel, error)
	GetQuoteByID(ctx context.Context, id uuid.UUID) (*QuotesModel, error)
	UpdateQuote(ctx context.Context, quoteID uuid.UUID, newQuote *QuoteRequest) (*QuotesModel, error)
	DeleteQuote(ctx context.Context, id uuid.UUID) error
}
