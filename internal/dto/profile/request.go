package profile

// UpdateUserRequest represents the user update request payload for updating user information.
type UpdateUserRequest struct {
	DisplayName string `json:"display_name" binding:"min=2,max=100"`
	Email       string `json:"email" binding:"email"`
}
