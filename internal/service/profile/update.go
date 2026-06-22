package profile

import (
	"context"
	"errors"

	"github.com/huypham67/bookmark-common/pkg/dbutils"
	profileDTO "github.com/huypham67/user-service/internal/dto/profile"
	"github.com/huypham67/user-service/internal/model"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

// UpdateUserInfo updates user's display name and email.
func (s *service) UpdateUserInfo(ctx context.Context, userID string, req profileDTO.UpdateUserRequest) error {
	segment := newrelic.FromContext(ctx).StartSegment("service.profile.UpdateUserInfo")
	defer segment.End()

	// Verify user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, dbutils.ErrRecordNotFoundType) {
			log.Warn().
				Str("user_id", userID).
				Msg("user not found")
		} else {
			log.Error().
				Err(err).
				Str("user_id", userID).
				Msg("failed to get user")
		}
		return err
	}

	// Check if the new email already exists for a different user (if email is provided)
	if req.Email != "" {
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil && !errors.Is(err, dbutils.ErrRecordNotFoundType) {
			log.Error().
				Err(err).
				Str("email", req.Email).
				Msg("failed to check if email exists")
			return err
		}

		// If email exists and belongs to a different user, return error
		if existingUser != nil && existingUser.ID != userID {
			log.Warn().
				Str("email", req.Email).
				Str("user_id", userID).
				Msg("email already registered to another user")
			return ErrEmailAlreadyRegistered
		}
	}

	user := &model.User{
		BaseModel: model.BaseModel{
			ID: userID,
		},
	}

	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to update user")
		return err
	}

	log.Info().
		Str("user_id", userID).
		Msg("user updated successfully")

	return nil
}
