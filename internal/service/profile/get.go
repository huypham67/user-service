package profile

import (
	"context"

	"github.com/huypham67/user-service/internal/model"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

// GetUserInfo retrieves user information by user ID.
func (s *service) GetUserInfo(ctx context.Context, userID string) (*model.User, error) {
	segment := newrelic.FromContext(ctx).StartSegment("service.profile.GetUserInfo")
	defer segment.End()

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
