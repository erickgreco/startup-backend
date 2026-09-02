package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LeadStatus string

const (
	LeadStatusNew       LeadStatus = "new"
	LeadStatusQuoted    LeadStatus = "quoted"
	LeadStatusConverted            = "converted"
	LeadStatusLost                 = "lost"
	LeadStatusMerged               = "merged"
)

type LeadsModel struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Message   string
	Status    LeadStatus //? new, quoted, converted, lost,	merged-en caso que exista un duplicado de cotizacion
	Source    LeadSource
	CreatedAt time.Time
	UpdatedAt time.Time
}

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

type LeadRepository interface {
	CreateLead(ctx context.Context, lead *LeadsModel) error
	GetLeadByStatus(ctx context.Context, status LeadStatus) (*LeadsModel, error)
}

type LeadService interface {
	CreateLead(ctx context.Context, lead *LeadsModel) error
}
