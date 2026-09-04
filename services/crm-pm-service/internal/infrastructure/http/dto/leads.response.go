package dto

import (
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/google/uuid"
)

type LeadResponse struct {
	ID           uuid.UUID         `json:"id"`
	Name         string            `json:"name"`
	Email        string            `json:"email"`
	Message      string            `json:"message"`
	DataQuality  domain.Quality    `json:"data_quality"`
	QualityScore float64           `json:"quality_score"`
	Status       domain.LeadStatus `json:"status"`
	Source       domain.LeadSource `json:"source"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func FromLeadDomain(l *domain.LeadsModel) LeadResponse {
	return LeadResponse{
		ID:           l.ID,
		Name:         l.Name,
		Email:        l.Email,
		Message:      l.Message,
		DataQuality:  l.DataQuality,
		QualityScore: l.QualityScore,
		Status:       l.Status,
		Source:       l.Source,
		CreatedAt:    l.CreatedAt,
		UpdatedAt:    l.UpdatedAt,
	}
}

func FromLeadDomainList(ls []*domain.LeadsModel) []LeadResponse {
	res := make([]LeadResponse, 0, len(ls))
	for _, l := range ls {
		res = append(res, FromLeadDomain(l))
	}
	return res
}
