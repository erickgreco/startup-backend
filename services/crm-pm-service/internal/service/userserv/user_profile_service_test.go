package userserv

import (
	"context"
	"errors"
	"testing"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
)

type mockUserProfileRepo struct {
	createFn func(ctx context.Context, up *domain.UserProfileModel) (*domain.UserProfileModel, error)
}

func (m *mockUserProfileRepo) CreateUserProfile(ctx context.Context, up *domain.UserProfileModel) (*domain.UserProfileModel, error) {
	return m.createFn(ctx, up)
}

func TestCreateUserProfile(t *testing.T) {
	tests := []struct {
		testName string
		req      *domain.UserProfileRequest
		repoFn   func(ctx context.Context, up *domain.UserProfileModel) (*domain.UserProfileModel, error)
		wantErr  error
	}{
		{
			testName: "Valid profile creation - Success",
			req: &domain.UserProfileRequest{
				FirstName:       "Erick",
				LastName:        "Greco",
				Email:           "test@test.com",
				PhoneNumber:     "+525512345678",
				Industry:        "Software Development",
				DefaultCurrency: "MXN",
				Language:        "es",
			},
			repoFn: func(ctx context.Context, up *domain.UserProfileModel) (*domain.UserProfileModel, error) {
				return up, nil
			},
			wantErr: nil,
		},
		{
			testName: "Invalid first name return error",
			req: &domain.UserProfileRequest{
				FirstName:       "   ",
				LastName:        "Greco",
				Email:           "test@test.com",
				PhoneNumber:     "+525512345678",
				Industry:        "Software Development",
				DefaultCurrency: "MXN",
				Language:        "es",
			},
			repoFn:  nil,
			wantErr: domain.ErrInvalidName,
		},
		{
			testName: "Invalid last name return error",
			req: &domain.UserProfileRequest{
				FirstName:       "Erick",
				LastName:        "    ",
				Email:           "test@test.com",
				PhoneNumber:     "+525512345678",
				Industry:        "Software Development",
				DefaultCurrency: "MXN",
				Language:        "es",
			},
			repoFn:  nil,
			wantErr: domain.ErrInvalidLastName,
		},
		{
			testName: "Invalid email return error",
			req: &domain.UserProfileRequest{
				FirstName:       "Erick",
				LastName:        "Greco",
				Email:           "test@test",
				PhoneNumber:     "+525512345678",
				Industry:        "Software Development",
				DefaultCurrency: "MXN",
				Language:        "es",
			},
			repoFn:  nil,
			wantErr: domain.ErrInvalidEmail,
		},
		{
			testName: "Invalid phone number return error",
			req: &domain.UserProfileRequest{
				FirstName:       "Erick",
				LastName:        "Greco",
				Email:           "test@test.com",
				PhoneNumber:     "551234567",
				Industry:        "Software Development",
				DefaultCurrency: "MXN",
				Language:        "es",
			},
			repoFn:  nil,
			wantErr: domain.ErrInvalidPhoneNumber,
		},
		{
			testName: "Invalid currency return error",
			req: &domain.UserProfileRequest{
				FirstName:       "Erick",
				LastName:        "Greco",
				Email:           "test@test.com",
				PhoneNumber:     "+525512345678",
				Industry:        "Software Development",
				DefaultCurrency: "XYZ",
				Language:        "es",
			},
			repoFn:  nil,
			wantErr: domain.ErrInvalidCurrency,
		},
		{
			testName: "Invalid language return error",
			req: &domain.UserProfileRequest{
				FirstName:       "Erick",
				LastName:        "Greco",
				Email:           "test@test.com",
				PhoneNumber:     "+525512345678",
				Industry:        "Software Development",
				DefaultCurrency: "MXN",
				Language:        "jp",
			},
			repoFn:  nil,
			wantErr: domain.ErrInvalidLanguage,
		},
		{
			testName: "Repo error propagates",
			req: &domain.UserProfileRequest{
				FirstName:       "Erick",
				LastName:        "Greco",
				Email:           "test@test.com",
				PhoneNumber:     "+525512345678",
				Industry:        "Software Development",
				DefaultCurrency: "MXN",
				Language:        "es",
			},
			repoFn: func(ctx context.Context, up *domain.UserProfileModel) (*domain.UserProfileModel, error) {
				return nil, domain.ErrInternal
			},
			wantErr: domain.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			repo := &mockUserProfileRepo{createFn: tt.repoFn}
			srv := NewUserProfileService(repo)

			result, err := srv.CreateUserProfile(context.Background(), tt.req)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				if result != nil {
					t.Fatal("expected nil result on error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("expected non-nil result")
			}
		})
	}
}
