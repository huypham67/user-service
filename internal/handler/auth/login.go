package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-common/pkg/requestutils"
	"github.com/huypham67/bookmark-common/pkg/response"
	authDTO "github.com/huypham67/user-service/internal/dto/auth"
	"github.com/huypham67/user-service/internal/service/auth"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

// Login handles the user login endpoint.
//
// @Summary Login User
// @Description Login a user with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body authDTO.LoginRequest true "User login data"
// @Success 200 {object} authDTO.LoginResponse "User logged in successfully"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 401 {object} gin.H "Invalid credentials"
// @Failure 404 {object} gin.H "User not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/users/login [post]
func (h *handler) Login(c *gin.Context) {
	segment := newrelic.FromContext(c).StartSegment("handler.auth.Login")
	defer segment.End()

	req, err := requestutils.Bind[authDTO.LoginRequest](c)

	if err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	token, err := h.service.LoginUser(c, *req)

	if err != nil {
		log.Error().
			Err(err).
			Str("username", req.Username).
			Msg("failed to login user")

		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			response.Unauthorized(c, "Invalid username or password")
		default:
			response.InternalServerError(c)
		}
		return
	}

	c.JSON(http.StatusOK, response.Success(token, "Logged in successfully!"))
}
