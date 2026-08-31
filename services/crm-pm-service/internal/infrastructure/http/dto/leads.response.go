package dto

import (
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/google/uuid"
)

type LeadResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"Updated_at"`
}

func FromLeadDomain(l *domain.LeadsModel) LeadResponse {
	return LeadResponse{
		ID:        l.ID,
		Name:      l.Name,
		Email:     l.Email,
		Message:   l.Message,
		Status:    l.Status,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}
}
