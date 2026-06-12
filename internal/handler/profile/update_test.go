package profile

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-common/pkg/jwt"
	profileDTO "github.com/huypham67/user-service/internal/dto/profile"
	"github.com/huypham67/user-service/internal/service/profile"
	"github.com/huypham67/user-service/internal/service/profile/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_UpdateUserInfo(t *testing.T) {
	t.Parallel()

	type expected struct {
		statusCode   int
		bodyContains string
	}

	testCases := []struct {
		name        string
		requestBody string
		setupClaims func(*gin.Context)
		setupMock   func(context.Context, *mocks.Service)
		expected    expected
	}{
		{
			name: "should return 200 when update succeeds",
			requestBody: `{
				"display_name":"Updated User",
				"email":"updated@example.com"
			}`,
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
						"UpdateUserInfo",
						ctx,
						"user-id",
						profileDTO.UpdateUserRequest{
							DisplayName: "Updated User",
							Email:       "updated@example.com",
						},
					).
					Return(nil).
					Once()
			},
			expected: expected{
				statusCode:   http.StatusOK,
				bodyContains: "Edit current user successfully!",
			},
		},
		{
			name:        "should return 400 when request body is invalid",
			requestBody: `{invalid json}`,
			setupClaims: func(ctx *gin.Context) {
				ctx.Set("claims", &jwt.CustomClaims{
					UserID: "user-id",
				})
			},
			setupMock: func(
				ctx context.Context,
				mockSvc *mocks.Service,
			) {
			},
			expected: expected{
				statusCode:   http.StatusBadRequest,
				bodyContains: "Invalid request body",
			},
		},
		{
			name: "should return 401 when claims are missing",
			requestBody: `{
				"display_name":"Updated User"
			}`,
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
			requestBody: `{
				"display_name":"Updated User"
			}`,
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
			name: "should return 409 when email already exists",
			requestBody: `{
				"display_name":"Updated User",
				"email":"existing@example.com"
			}`,
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
						"UpdateUserInfo",
						ctx,
						"user-id",
						profileDTO.UpdateUserRequest{
							DisplayName: "Updated User",
							Email:       "existing@example.com",
						},
					).
					Return(
						profile.ErrEmailAlreadyRegistered,
					).
					Once()
			},
			expected: expected{
				statusCode:   http.StatusConflict,
				bodyContains: "Email already exists",
			},
		},
		{
			name: "should return 500 when service fails",
			requestBody: `{
				"display_name":"Updated User",
				"email":"updated@example.com"
			}`,
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
						"UpdateUserInfo",
						ctx,
						"user-id",
						profileDTO.UpdateUserRequest{
							DisplayName: "Updated User",
							Email:       "updated@example.com",
						},
					).
					Return(
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

			req := httptest.NewRequest(
				http.MethodPut,
				"/v1/self/info",
				bytes.NewBufferString(tc.requestBody),
			)

			req.Header.Set("Content-Type", "application/json")

			ctx.Request = req

			tc.setupClaims(ctx)
			tc.setupMock(ctx, mockSvc)

			handler := NewHandler(mockSvc)
			handler.UpdateUserInfo(ctx)

			assert.Equal(t, tc.expected.statusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expected.bodyContains)

			mockSvc.AssertExpectations(t)
		})
	}
}
