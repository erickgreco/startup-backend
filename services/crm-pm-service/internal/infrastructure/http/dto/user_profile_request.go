package dto

import "github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"

type UserProfileRequest struct {
	FirstName       string  `json:"first_name" binding:"required"`
	LastName        string  `json:"last_name" binding:"required"`
	Email           string  `json:"email" binding:"required"`
	PhoneNumber     string  `json:"phone_number" binding:"required"`
	ProfileImageURL *string `json:"profile_image_url,omitempty"`
	BusinessName    *string `json:"business_name,omitempty"`
	Industry        string  `json:"industry" binding:"required"`
	LogoURL         *string `json:"logo_url,omitempty"`
	DefaultCurrency string  `json:"default_currency" binding:"required"`
	Language        string  `json:"language" binding:"required"`
	WebsiteURL      *string `json:"website_url,omitempty"`
}

func (r UserProfileRequest) ToDomain() *domain.UserProfileRequest {
	return &domain.UserProfileRequest{
		FirstName:       r.FirstName,
		LastName:        r.LastName,
		Email:           r.Email,
		PhoneNumber:     r.PhoneNumber,
		ProfileImageURL: r.ProfileImageURL,
		BusinessName:    r.BusinessName,
		Industry:        r.Industry,
		LogoURL:         r.LogoURL,
		DefaultCurrency: r.DefaultCurrency,
		Language:        r.Language,
		WebsiteURL:      r.WebsiteURL,
	}
}
