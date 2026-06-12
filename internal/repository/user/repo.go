package user

import (
	"context"

	"github.com/huypham67/user-service/internal/model"
	"gorm.io/gorm"
)

// Repository defines the interface for user repository, which provides methods to interact with the user data in the database.
//
//go:generate mockery --name=Repository --output=./mocks --outpkg=mocks --filename=mock_repo.go
type Repository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository with the given GORM database client.
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}
