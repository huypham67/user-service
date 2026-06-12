package profile

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-common/pkg/jwt"
	"github.com/huypham67/bookmark-common/pkg/requestutils"
	"github.com/huypham67/bookmark-common/pkg/response"
	profileDTO "github.com/huypham67/user-service/internal/dto/profile"
	"github.com/huypham67/user-service/internal/service/profile"
	"github.com/rs/zerolog/log"
)

// UpdateUserInfo handles the user info update endpoint.
//
// @Summary Update User Info
// @Description Update authenticated user's display name and email
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body profileDTO.UpdateUserRequest true "User update data"
// @Success 200 {object} profileDTO.UpdateUserResponse "User updated successfully"
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 409 {object} gin.H "Email already exists"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /v1/self/info [put]
func (h *handler) UpdateUserInfo(c *gin.Context) {
	userID, err := jwt.GetUserIDFromContext(c)

	if err != nil {
		log.Warn().Msg("user ID not found in context")
		response.Unauthorized(c, "Unauthorized")
		return
	}

	req, err := requestutils.Bind[profileDTO.UpdateUserRequest](c)

	if err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.service.UpdateUserInfo(c, userID, *req); err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Str("email", req.Email).
			Msg("failed to update user info")

		if errors.Is(err, profile.ErrEmailAlreadyRegistered) {
			response.Conflict(c, "Email already exists")
			return
		}

		response.InternalServerError(c)

		return
	}

	c.JSON(http.StatusOK, response.Message("Edit current user successfully!"))
}
