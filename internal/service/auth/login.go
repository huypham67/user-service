package auth

import (
	"context"
	"errors"

	"github.com/huypham67/bookmark-common/pkg/dbutils"
	"github.com/huypham67/bookmark-common/pkg/security"
	authDTO "github.com/huypham67/user-service/internal/dto/auth"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

const (
	eventLoginFailure  = "LoginFailure"
	metricLoginFailure = "Custom/Auth/LoginFailure"
	reasonInvalidCreds = "invalid_credentials"
)

func recordLoginFailure(ctx context.Context) {
	app := newrelic.FromContext(ctx).Application()
	app.RecordCustomMetric(metricLoginFailure, 1)
	app.RecordCustomEvent(eventLoginFailure, map[string]interface{}{
		"reason": reasonInvalidCreds,
	})
}

// LoginUser authenticates a user by validating credentials and returns a JWT token.
func (s *service) LoginUser(ctx context.Context, req authDTO.LoginRequest) (string, error) {
	segment := newrelic.FromContext(ctx).StartSegment("service.auth.LoginUser")
	defer segment.End()

	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, dbutils.ErrRecordNotFoundType) {
			recordLoginFailure(ctx)
			return "", ErrInvalidCredentials
		}
		log.Error().
			Err(err).
			Str("username", req.Username).
			Msg("failed to get user by username")
		return "", ErrInternalServerError
	}

	if user == nil {
		recordLoginFailure(ctx)
		return "", ErrInvalidCredentials
	}

	compareErr := s.passwordHasher.Compare(user.Password, req.Password)
	if err := compareErr; err != nil {
		if errors.Is(err, security.ErrPasswordMismatch) {
			recordLoginFailure(ctx)
			return "", ErrInvalidCredentials
		}
		log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to compare password hash")
		return "", ErrInternalServerError
	}

	token, err := s.tokenGenerator.GenerateToken(user.ID, user.DisplayName, user.Email)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to generate token")
		return "", ErrInternalServerError
	}

	log.Info().
		Str("user_id", user.ID).
		Str("username", user.Username).
		Msg("user logged in successfully")

	return token, nil
}
