package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserProfileRequest struct {
	FirstName       string
	LastName        string
	Email           string
	PhoneNumber     string
	ProfileImageURL *string
	BusinessName    *string
	Industry        string
	LogoURL         *string
	DefaultCurrency string
	Language        string
	WebsiteURL      *string
}

type UserProfileModel struct {
	ID              uuid.UUID
	FirstName       string
	LastName        string
	Email           string
	PhoneNumber     string
	ProfileImageURL *string
	BusinessName    *string
	Industry        string
	LogoURL         *string
	DefaultCurrency Currency
	Timezone        string
	Language        string
	WebsiteURL      *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UserProfileRepository interface {
	CreateUserProfile(ctx context.Context, up *UserProfileModel) (*UserProfileModel, error)
}

type UserProfileService interface {
	CreateUserProfile(ctx context.Context, req *UserProfileRequest) (*UserProfileModel, error)
}
