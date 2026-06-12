package profile

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-common/pkg/jwt"
	"github.com/huypham67/user-service/internal/model"
	"github.com/huypham67/user-service/internal/service/profile/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetUserInfo(t *testing.T) {
	t.Parallel()

	type expected struct {
		statusCode   int
		bodyContains string
	}

	testCases := []struct {
		name        string
		setupClaims func(*gin.Context)
		setupMock   func(context.Context, *mocks.Service)
		expected    expected
	}{
		{
			name: "should return 200 when user exists",
			setupClaims: func(ctx *gin.Context) {
				ctx.Set("claims", &jwt.CustomClaims{
					UserID: "user-id",
				})
			},
			setupMock: func(
				ctx context.Context,
				mockSvc *mocks.Service,
			) {
				mockSvc.
					On(
						"GetUserInfo",
						ctx,
						"user-id",
					).
					Return(
						&model.User{
							BaseModel: model.BaseModel{
								ID: "user-id",
							},
							DisplayName: "John Doe",
							Username:    "johndoe",
							Email:       "john@example.com",
						},
						nil,
					).
					Once()
			},
			expected: expected{
				statusCode:   http.StatusOK,
				bodyContains: "User information retrieved successfully!",
			},
		},
		{
			name: "should return 401 when claims are missing",
			setupClaims: func(ctx *gin.Context) {
			},
			setupMock: func(
				ctx context.Context,
				mockSvc *mocks.Service,
			) {
			},
			expected: expected{
				statusCode:   http.StatusUnauthorized,
				bodyContains: "Unauthorized",
			},
		},
		{
			name: "should return 401 when claims type is invalid",
			setupClaims: func(ctx *gin.Context) {
				ctx.Set(
					"claims",
					"invalid-claims",
				)
			},
			setupMock: func(
				ctx context.Context,
				mockSvc *mocks.Service,
			) {
			},
			expected: expected{
				statusCode:   http.StatusUnauthorized,
				bodyContains: "Unauthorized",
			},
		},
		{
			name: "should return 500 when service fails",
			setupClaims: func(ctx *gin.Context) {
				ctx.Set("claims", &jwt.CustomClaims{
					UserID: "user-id",
				})
			},
			setupMock: func(
				ctx context.Context,
				mockSvc *mocks.Service,
			) {
				mockSvc.
					On(
						"GetUserInfo",
						ctx,
						"user-id",
					).
					Return(
						nil,
						errors.New("database error"),
					).
					Once()
			},
			expected: expected{
				statusCode:   http.StatusInternalServerError,
				bodyContains: "Internal Server Error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gin.SetMode(gin.TestMode)

			mockSvc := new(mocks.Service)

			recorder := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(recorder)

			req := httptest.NewRequest("GET", "/v1/self/info", nil)

			ctx.Request = req

			tc.setupClaims(ctx)
			tc.setupMock(ctx, mockSvc)

			handler := NewHandler(mockSvc)
			handler.GetUserInfo(ctx)

			assert.Equal(t, tc.expected.statusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expected.bodyContains)

			mockSvc.AssertExpectations(t)
		})
	}
}
