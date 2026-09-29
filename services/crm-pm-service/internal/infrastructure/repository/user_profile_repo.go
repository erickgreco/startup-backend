package repository

import (
	"context"
	"fmt"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/repository/sqlgen"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserProfileRepository struct {
	queries *sqlgen.Queries
}

func NewUserProfileRepository(pool *pgxpool.Pool) *UserProfileRepository {
	return &UserProfileRepository{
		queries: sqlgen.New(pool),
	}
}

func (r *UserProfileRepository) CreateUserProfile(ctx context.Context, up *domain.UserProfileModel) (*domain.UserProfileModel, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeDuration)
	defer cancel()

	row, err := r.queries.CreateUserProfile(ctx, userProfileDomainToGen(up))
	if err != nil {
		return nil, fmt.Errorf("user_profile_repo: create_user_profile: %w", translateWriteError(err))
	}

	created, err := genToUserProfileDomain(row)
	if err != nil {
		return nil, fmt.Errorf("user_profile_repo: create_user_profile: map_result: %w", fmt.Errorf("%w: %w", domain.ErrInternal, err))
	}

	return created, nil
}
