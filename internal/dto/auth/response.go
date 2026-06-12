package auth

import (
	"time"

	"github.com/huypham67/bookmark-common/pkg/response"
)

// UserData represents the user data in the response.
type UserData struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
}

// RegisterUserResponse is a type alias for user registration response.
type RegisterUserResponse = response.SuccessResponse[UserData]

// LoginResponse is a type alias for login response with JWT token.
type LoginResponse = response.SuccessResponse[string]
