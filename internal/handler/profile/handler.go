package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/huypham67/user-service/internal/service/profile"
)

// Handler defines the interface for profile-related HTTP handlers, including endpoints for retrieving and updating user information.
type Handler interface {
	GetUserInfo(c *gin.Context)
	UpdateUserInfo(c *gin.Context)
}

type handler struct {
	service profile.Service
}

// NewHandler creates a new instance of the profile handler with the provided profile service.
func NewHandler(service profile.Service) Handler {
	return &handler{
		service: service,
	}
}
