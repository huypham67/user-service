package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-common/pkg/jwt"
	"github.com/huypham67/bookmark-common/pkg/response"
	profileDTO "github.com/huypham67/user-service/internal/dto/profile"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
)

// GetUserInfo handles the user info endpoint.
//
// @Summary Get User Info
// @Description Get authenticated user information from JWT token
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} profileDTO.UserResponse "User information"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "User not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/self/info [get]
func (h *handler) GetUserInfo(c *gin.Context) {
	segment := newrelic.FromContext(c).StartSegment("handler.profile.GetUserInfo")
	defer segment.End()

	userID, err := jwt.GetUserIDFromContext(c)

	if err != nil {
		log.Warn().Msg("user ID not found in context")
		response.Unauthorized(c, "Unauthorized")
		return
	}

	user, err := h.service.GetUserInfo(c, userID)

	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to get user info")

		response.InternalServerError(c)

		return
	}

	c.JSON(http.StatusOK, response.Success(
		profileDTO.UserData{
			ID:          user.ID,
			DisplayName: user.DisplayName,
			Username:    user.Username,
			Email:       user.Email,
			CreatedAt:   user.CreatedAt,
		},
		"User information retrieved successfully!",
	))
}
