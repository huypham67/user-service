package profile

import (
	"context"

	"github.com/huypham67/user-service/internal/model"
	"github.com/rs/zerolog/log"
)

// GetUserInfo retrieves user information by user ID.
func (s *service) GetUserInfo(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to get user by ID")
		return nil, err
	}

	return user, nil
}
