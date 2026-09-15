package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-playground/locales/currency"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type QuoteRequest struct {
	ClientName                  string
	ClientBusiness              *string
	ClientEmail                 string
	ClientPhone                 string
	ProjectTitle                string
	ProjectDescription          string
	ProjectEstimatedDurationSec time.Duration
	Items                       []QuoteItem
	PaymentMethod               PayMethod
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
	ProjectTitle                string
	ProjectDescription          string
	ProjectResume               *string
	ProjectEstimatedDurationSec time.Duration

	// Detalles de la oferta
	FunctionalityExample *string
	Items                []QuoteItem // ** deliverable
	PaymentMethod        PayMethod

	// Estados y fechas
	Status     QuoteStatus
	CreatedAt  time.Time
	ValidUntil time.Time // ** cuenta a partir de cuando se envia
	SentAt     *time.Time
	AcceptedAt *time.Time
	RejectedAt *time.Time

	// Totales
	Subtotal Money
	Total    Money

	// Metadata
	Currency currency.Type
	Notes    *string

	// Auditoria
	UpdatedAt time.Time
}

// ** Helper: serializar Items a JSONB (para insertar en DB)
func (q *QuotesModel) SetItemsToJSON() ([]byte, error) {
	return json.Marshal(q.Items)
}

// ** Helper: deserializar JSONB a Items (para leer de BD)
func (q *QuotesModel) SetItemsFromJSON(data []byte) error {
	return json.Unmarshal(data, &q.Items)
}

// ** Representa dinero en centavos 1050 = $10.50 USD
type Money = decimal.Decimal

type PayMethod string

const (
	/*
		- Proyectos pequenos
		- Clientes nuevos sin track record
		- Trabajos de bajo riesgo (plantillas o disenos simples)
	*/
	PaymentUpfront = "upfront"
	/*
		- Desarrollo complejo
		- Clientes corporativos
		- Relaciones a largo plazo
	*/
	PaymentMilestone = "milestone"
	/*
		- Proyectos pequenos
		- Clientes con historial establecido
		- Entregables simples/tangibles
		- NO RECOMENDADO PARA PROYECTOS GRANDES
	*/
	PaymentUponDelivery = "upon_delivery"
	/*
		- Proyectos costosos
		- Clientes nuevos pero con referencias
		- Proteccion sin ahuyentar al cliente
	*/
	PaymentSplit = "split"
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
	QuoteStatusDraft     QuoteStatus = "new"
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
}
