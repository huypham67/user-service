package user

import (
	"context"

	"github.com/huypham67/bookmark-common/pkg/dbutils"
	"github.com/huypham67/user-service/internal/model"
)

// GetByEmail retrieves a user by their email address.
func (r *repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.getUserByField(ctx, "email", email)
}

// GetByUsername retrieves a user by their username.
func (r *repository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.getUserByField(ctx, "username", username)
}

// GetByID retrieves a user by their ID.
func (r *repository) GetByID(ctx context.Context, userID string) (*model.User, error) {
	return r.getUserByField(ctx, "id", userID)
}

func (r *repository) getUserByField(ctx context.Context, fieldName, fieldValue string) (*model.User, error) {
	var user *model.User
	if err := r.db.WithContext(ctx).Where(fieldName+" = ?", fieldValue).First(&user).Error; err != nil {
		return nil, dbutils.ClassifyError(err)
	}
	return user, nil
}
