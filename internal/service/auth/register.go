package auth

import (
	"context"
	"errors"

	"github.com/huypham67/bookmark-common/pkg/dbutils"
	authDTO "github.com/huypham67/user-service/internal/dto/auth"
	"github.com/huypham67/user-service/internal/model"
	"github.com/rs/zerolog/log"
)

// RegisterUser registers a new user by hashing the password and saving to the database.
func (s *service) RegisterUser(ctx context.Context, req authDTO.RegisterUserRequest) (*model.User, error) {
	hashedPassword, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to hash password")
		return nil, ErrInternalServerError
	}

	user := &model.User{
		DisplayName: req.DisplayName,
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, dbutils.ErrDuplicationType) {
			log.Warn().
				Str("email", req.Email).
				Str("username", req.Username).
				Msg("user already exists")
			return nil, ErrUserAlreadyExists
		}

		log.Error().
			Err(err).
			Str("email", req.Email).
			Str("username", req.Username).
			Msg("failed to register user")
		return nil, ErrInternalServerError
	}

	log.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Str("username", user.Username).
		Msg("user registered successfully")

	return user, nil
}
