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

// Register handles the user registration endpoint.
//
// @Summary Register User
// @Description Register a new user with email, username, and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body authDTO.RegisterUserRequest true "User registration data"
// @Success 201 {object} authDTO.RegisterUserResponse "User registered successfully"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 409 {object} gin.H "User already exists"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/users/register [post]
func (h *handler) Register(c *gin.Context) {
	segment := newrelic.FromContext(c).StartSegment("handler.auth.Register")
	defer segment.End()

	req, err := requestutils.Bind[authDTO.RegisterUserRequest](c)

	if err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	user, err := h.service.RegisterUser(c, *req)

	if err != nil {
		log.Error().
			Err(err).
			Str("email", req.Email).
			Str("username", req.Username).
			Msg("failed to register user")

		switch {
		case errors.Is(err, auth.ErrUserAlreadyExists):
			response.Conflict(c, "User already exists")
		default:
			response.InternalServerError(c)
		}
		return
	}

	c.JSON(http.StatusCreated, response.Success(
		authDTO.UserData{
			ID:          user.ID,
			DisplayName: user.DisplayName,
			Username:    user.Username,
			Email:       user.Email,
			CreatedAt:   user.CreatedAt,
		},
		"Register an user successfully!",
	))
}
