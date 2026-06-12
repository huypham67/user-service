package auth

import (
	"context"
	"errors"

	"github.com/huypham67/bookmark-common/pkg/dbutils"
	"github.com/huypham67/bookmark-common/pkg/security"
	authDTO "github.com/huypham67/user-service/internal/dto/auth"
	"github.com/rs/zerolog/log"
)

// LoginUser authenticates a user by validating credentials and returns a JWT token.
func (s *service) LoginUser(ctx context.Context, req authDTO.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, dbutils.ErrRecordNotFoundType) {
			return "", ErrInvalidCredentials
		}
		log.Error().
			Err(err).
			Str("username", req.Username).
			Msg("failed to get user by username")
		return "", ErrInternalServerError
	}

	if user == nil {
		return "", ErrInvalidCredentials
	}

	if err := s.passwordHasher.Compare(user.Password, req.Password); err != nil {
		if errors.Is(err, security.ErrPasswordMismatch) {
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
