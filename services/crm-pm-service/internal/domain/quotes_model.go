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
	ProjectEstimatedDurationDays int
	Items                        []ItemRequest
	PaymentMethod                PayMethod
	Currency                     Currency // TODO: trabajar con location una vez se implemente middleware
	Notes                        *string
}

type QuotesModel struct {
	ID          uuid.UUID
	LeadID      uuid.UUID
	QuoteNumber string

	// Client data
	ClientName     string
	ClientBusiness *string
	ClientEmail    string
	ClientPhone    string

	// Project data
	ProjectTitle                 string
	ProjectDescription           string
	ProjectResume                *string
	ProjectEstimatedDurationDays int

	// Detalles de la oferta
	FunctionalityExample *string
	Items                []QuoteItem // ** deliverable
	PaymentMethod        PayMethod

	// Estados y fechas
	Status     QuoteStatus
	CreatedAt  time.Time
	ValidUntil *time.Time // ** cuenta a partir de cuando se envia
	SentAt     *time.Time
	AcceptedAt *time.Time
	RejectedAt *time.Time

	// Totales
	Subtotal Money
	Total    Money

	// Metadata
	Currency Currency
	Notes    *string

	// Auditoria
	UpdatedAt time.Time
}

type PayMethod string

const (
	/*
		- Proyectos pequenos
		- Clientes nuevos sin track record
		- Trabajos de bajo riesgo (plantillas o disenos simples)
	*/
	PaymentUpfront PayMethod = "upfront"
	/*
		- Desarrollo complejo
		- Clientes corporativos
		- Relaciones a largo plazo
	*/
	PaymentMilestone PayMethod = "milestone"
	/*
		- Proyectos pequenos
		- Clientes con historial establecido
		- Entregables simples/tangibles
		- NO RECOMENDADO PARA PROYECTOS GRANDES
	*/
	PaymentUponDelivery PayMethod = "upon_delivery"
	/*
		- Proyectos costosos
		- Clientes nuevos pero con referencias
		- Proteccion sin ahuyentar al cliente
	*/
	PaymentSplit PayMethod = "split"
)

func (p PayMethod) IsValid() bool {
	switch p {
	case PaymentUpfront, PaymentMilestone, PaymentUponDelivery, PaymentSplit:
		return true
	default:
		return false
	}
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
}

type QuoteService interface {
	CreateQuote(ctx context.Context, leadID uuid.UUID, req *QuoteRequest) (*QuotesModel, error)
	GetQuoteByID(ctx context.Context, id uuid.UUID) (*QuotesModel, error)
}
