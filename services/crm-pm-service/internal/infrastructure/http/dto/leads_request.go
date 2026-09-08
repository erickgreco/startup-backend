package dto

import "github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"

type CreateLeadRequest struct {
	Name         string `json:"name" binding:"required"`
	BusinessName string `json:"business_name" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Message      string `json:"message" binding:"required"`
}

func (r CreateLeadRequest) ToDomain() *domain.LeadRequest {
	return &domain.LeadRequest{
		Name:         r.Name,
		BusinessName: r.BusinessName,
		PhoneNumber:  r.PhoneNumber,
		Email:        r.Email,
		Message:      r.Message,
	}
}

type UpdateLeadStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=new quoted converted lost merged"`
}
