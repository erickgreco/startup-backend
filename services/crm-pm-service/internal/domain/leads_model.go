package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LeadsModel struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Message   string
	Status    string //? new, quoted, converted, lost
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LeadRepository interface {
	CreateLead(ctx context.Context, lead *LeadsModel) error
}

type LeadService interface {
	CreateLead(ctx context.Context, lead *LeadsModel) error
}
