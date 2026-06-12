package profile

import (
	"time"

	"github.com/huypham67/bookmark-common/pkg/response"
)

// UserData represents the user data in the profile response.
type UserData struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserResponse is a type alias for user profile response.
type UserResponse = response.SuccessResponse[UserData]

// UpdateUserResponse is a type alias for user update response.
type UpdateUserResponse = response.SuccessResponse[struct{}]
