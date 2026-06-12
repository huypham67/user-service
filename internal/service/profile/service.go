// Package profile provides user profile services.
// It handles user profile retrieval and update operations.
package profile

import (
	"context"
	"errors"

	profileDTO "github.com/huypham67/user-service/internal/dto/profile"
	"github.com/huypham67/user-service/internal/model"
	"github.com/huypham67/user-service/internal/repository/user"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)

// Service defines the contract for profile operations.
//
//go:generate mockery --name=Service --output=./mocks --outpkg=mocks --filename=mock_service.go
type Service interface {
	GetUserInfo(ctx context.Context, userID string) (*model.User, error)
	UpdateUserInfo(ctx context.Context, userID string, req profileDTO.UpdateUserRequest) error
}

type service struct {
	userRepo user.Repository
}

// NewService creates a new profile service with the provided user repository.
func NewService(userRepo user.Repository) Service {
	return &service{
		userRepo: userRepo,
	}
}
