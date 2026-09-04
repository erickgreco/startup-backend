package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LeadsModel struct {
	ID           uuid.UUID
	Name         string
	Email        string
	Message      string
	DataQuality  Quality
	QualityScore float64    // 0-100
	Status       LeadStatus //** new, quoted, converted, lost,	merged-en caso que exista un duplicado de cotizacion
	Source       LeadSource
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Quality string

const (
	High   Quality = "high"   //** Info clara, scope definido
	Medium Quality = "medium" //** Mas del 50%, aun con dudas
	Low    Quality = "low"    //** Falta contexto
)

// ** Breakdown de analisis
type QualityAnalisis struct {
	Score         float64
	QA            Quality
	Indicators    map[string]float64 //** Debug info
	FoundConcepts []string           // ** Conceptos detectados
}

type LeadStatus string

const (
	LeadStatusNew       LeadStatus = "new"
	LeadStatusQuoted    LeadStatus = "quoted"
	LeadStatusConverted LeadStatus = "converted"
	LeadStatusLost      LeadStatus = "lost"
	LeadStatusMerged    LeadStatus = "merged"
)

type LeadSource string

const (
	LeadSrcLandingPage LeadSource = "landing_page"
	LeadSrcEmail       LeadSource = "email"
	LeadSrcWsApp       LeadSource = "whatsapp"
	LeadSrcManual      LeadSource = "manual"
)

func (s LeadSource) IsValid() bool {
	switch s {
	case LeadSrcLandingPage, LeadSrcEmail, LeadSrcWsApp, LeadSrcManual:
		return true
	default:
		return false
	}
}

func (s LeadStatus) IsValid() bool {
	switch s {
	case LeadStatusNew, LeadStatusQuoted, LeadStatusConverted, LeadStatusLost, LeadStatusMerged:
		return true
	default:
		return false
	}
}

type LeadRepository interface {
	CreateLead(ctx context.Context, lead *LeadsModel) error
	GetLeadsByStatus(ctx context.Context, status LeadStatus) ([]*LeadsModel, error)
	GetLeadByID(ctx context.Context, id uuid.UUID) (*LeadsModel, error)
	UpdateLeadStatus(ctx context.Context, id uuid.UUID, status LeadStatus) (*LeadsModel, error)
	DeleteLead(ctx context.Context, id uuid.UUID) error
}

type LeadService interface {
	CreateLead(ctx context.Context, lead *LeadsModel) error
	GetLeadsByStatus(ctx context.Context, status LeadStatus) ([]*LeadsModel, error)
	GetLeadByID(ctx context.Context, id uuid.UUID) (*LeadsModel, error)
	UpdateLeadStatus(ctx context.Context, id uuid.UUID, status LeadStatus) (*LeadsModel, error)
	DeleteLead(ctx context.Context, id uuid.UUID) error
}
