package userserv

import (
	"context"
	"fmt"
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/shared/idgen"
)

type UserProfileService struct {
	repo domain.UserProfileRepository
}

func NewUserProfileService(repo domain.UserProfileRepository) *UserProfileService {
	return &UserProfileService{
		repo: repo,
	}
}

// TODO: corregir timezone "Local"
func (s *UserProfileService) CreateUserProfile(ctx context.Context, req *domain.UserProfileRequest) (*domain.UserProfileModel, error) {
	name, lastName, err := validateFullName(req.FirstName, req.LastName)
	if err != nil {
		return nil, fmt.Errorf("user_service: create_user_profile: %w", err)
	}

	validatedEmail, err := validateEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("user_service: create_user_profile: %w", err)
	}

	validatedPhoneNumber, err := validatePhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("user_service: create_user_profile: %w", err)
	}

	cur, err := domain.ParseCurrency(req.DefaultCurrency)
	if err != nil {
		return nil, fmt.Errorf("user_service: create_user_profile: %w", err)
	}

	timeZone, err := time.LoadLocation("Local")
	if err != nil {
		return nil, fmt.Errorf("user_service: create_user_profile: %w", err)
	}

	if !isValidLanguage(req.Language) {
		return nil, fmt.Errorf("user_service: create_user:profile: %w", domain.ErrInvalidLanguage)
	}

	now := time.Now()

	profile := &domain.UserProfileModel{
		ID:              idgen.NewID(),
		FirstName:       name,
		LastName:        lastName,
		Email:           validatedEmail,
		PhoneNumber:     validatedPhoneNumber,
		ProfileImageURL: req.ProfileImageURL,
		BusinessName:    req.BusinessName,
		Industry:        req.Industry,
		LogoURL:         req.LogoURL,
		DefaultCurrency: cur,
		Timezone:        timeZone.String(),
		Language:        req.Language,
		WebsiteURL:      req.WebsiteURL,
		CreatedAt:       now,
	}

	created, err := s.repo.CreateUserProfile(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("user_service: create_user_profile: %w", err)
	}

	return created, nil
}
