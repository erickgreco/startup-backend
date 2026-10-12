package dto

import (
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/google/uuid"
)

type LeadResponseUser struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	BusinessName string    `json:"business_name"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Message      string    `json:"message"`
	DataQuality  string    `json:"data_quality"`
	QualityScore float64   `json:"quality_score"`
	Status       string    `json:"status"`
	Source       string    `json:"source"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LeadResponseClient struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	BusinessName string    `json:"business_name"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func FromLeadDomain(l *domain.LeadsModel) any {
	switch l.Source {
	case domain.LeadSrcManual:
		return toLeadResponseUser(l)
	default:
		return toLeadResponseClient(l)
	}
}

func toLeadResponseUser(l *domain.LeadsModel) LeadResponseUser {
	return LeadResponseUser{
		ID:           l.ID,
		Name:         l.Name,
		BusinessName: l.BusinessName,
		PhoneNumber:  l.PhoneNumber,
		Email:        l.Email,
		Message:      l.Message,
		DataQuality:  string(l.DataQuality),
		QualityScore: l.QualityScore,
		Status:       string(l.Status),
		Source:       string(l.Source),
		CreatedAt:    l.CreatedAt,
		UpdatedAt:    l.UpdatedAt,
	}
}

func toLeadResponseClient(l *domain.LeadsModel) LeadResponseClient {
	return LeadResponseClient{
		ID:           l.ID,
		Name:         l.Name,
		BusinessName: l.BusinessName,
		PhoneNumber:  l.PhoneNumber,
		Email:        l.Email,
		Message:      l.Message,
		CreatedAt:    l.CreatedAt,
		UpdatedAt:    l.UpdatedAt,
	}
}

func FromLeadDomainList(ls []*domain.LeadsModel) []any {
	res := make([]any, 0, len(ls))
	for _, l := range ls {
		res = append(res, FromLeadDomain(l))
	}
	return res
}
