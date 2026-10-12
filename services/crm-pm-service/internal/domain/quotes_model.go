package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PayMethod string

const (
	PaymentUpfront      PayMethod = "upfront"
	PaymentMilestone    PayMethod = "milestone"
	PaymentUponDelivery PayMethod = "upon_delivery"
	PaymentSplit        PayMethod = "split"
)

type QuoteStatus string

const (
	QuoteStatusDraft     QuoteStatus = "draft"
	QuoteStatusSent      QuoteStatus = "sent"
	QuoteStatusAccepted  QuoteStatus = "accepted"
	QuoteStatusRejected  QuoteStatus = "rejected"
	QuoteStatusExpired   QuoteStatus = "expired"
	QuoteStatusConverted QuoteStatus = "converted"
)

const DefaultQuoteValidityDays = 30

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

func (p PayMethod) IsValid() bool {
	switch p {
	case PaymentUpfront, PaymentMilestone, PaymentUponDelivery, PaymentSplit:
		return true
	default:
		return false
	}
}

func (s QuoteStatus) IsValid() bool {
	switch s {
	case QuoteStatusDraft, QuoteStatusSent, QuoteStatusAccepted, QuoteStatusRejected, QuoteStatusExpired, QuoteStatusConverted:
		return true
	default:
		return false
	}
}

func (s QuoteStatus) CanBeSent() bool {
	switch s {
	case QuoteStatusDraft:
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
	SendQuote(ctx context.Context, quoteID uuid.UUID, sentAt, validUntil time.Time) error
}

type QuoteService interface {
	CreateManualQuote(ctx context.Context, req *QuoteRequest) (*QuotesModel, error)
	CreateQuoteFromLead(ctx context.Context, leadID uuid.UUID, req *QuoteRequest) (*QuotesModel, error)
	GetQuoteByID(ctx context.Context, id uuid.UUID) (*QuotesModel, error)
	UpdateQuote(ctx context.Context, quoteID uuid.UUID, newQuote *QuoteRequest) (*QuotesModel, error)
	DeleteQuote(ctx context.Context, id uuid.UUID) error
	SendQuote(ctx context.Context, quoteID uuid.UUID) error
}
