package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/huypham67/user-service/internal/service/auth"
)

// Handler defines the interface for authentication-related HTTP handlers, including user registration and login.
type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type handler struct {
	service auth.Service
}

// NewHandler creates a new instance of the auth handler with the provided authentication service.
func NewHandler(service auth.Service) Handler {
	return &handler{
		service: service,
	}
}
