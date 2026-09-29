package dto

import (
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
)

type UserProfileResponse struct {
	ID              string    `json:"id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	PhoneNumber     string    `json:"phone_number"`
	ProfileImageURL *string   `json:"profile_image_url"`
	BusinessName    *string   `json:"business_name"`
	Industry        string    `json:"industry"`
	LogoURL         *string   `json:"logo_url"`
	DefaultCurrency string    `json:"default_currency"`
	Timezone        string    `json:"timezone"`
	Language        string    `json:"language"`
	WebsiteURL      *string   `json:"website_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func FromUserProfileDomain(up *domain.UserProfileModel) UserProfileResponse {
	return UserProfileResponse{
		ID:              up.ID.String(),
		FirstName:       up.FirstName,
		LastName:        up.LastName,
		Email:           up.Email,
		PhoneNumber:     up.PhoneNumber,
		ProfileImageURL: up.ProfileImageURL,
		BusinessName:    up.BusinessName,
		Industry:        up.Industry,
		LogoURL:         up.LogoURL,
		DefaultCurrency: up.DefaultCurrency.String(),
		Timezone:        up.Timezone,
		Language:        up.Language,
		WebsiteURL:      up.WebsiteURL,
		CreatedAt:       up.CreatedAt,
		UpdatedAt:       up.UpdatedAt,
	}
}
