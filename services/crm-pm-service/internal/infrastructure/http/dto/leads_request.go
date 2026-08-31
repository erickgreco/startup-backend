package dto

import "github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"

type CreateLeadRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Message string `json:"message" binding:"required"`
}

func (r CreateLeadRequest) ToDomain() *domain.LeadsModel {
	return &domain.LeadsModel{
		Name:    r.Name,
		Email:   r.Email,
		Message: r.Message,
	}
}
