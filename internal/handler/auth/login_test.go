package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	authDTO "github.com/huypham67/user-service/internal/dto/auth"
	"github.com/huypham67/user-service/internal/service/auth"
	"github.com/huypham67/user-service/internal/service/auth/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Login(t *testing.T) {
	t.Parallel()

	type expected struct {
		statusCode   int
		bodyContains string
	}

	testCases := []struct {
		name        string
		requestBody string
		setupMock   func(context.Context, *mocks.Service)
		expected    expected
	}{
		{
			name: "should return 200 when login succeeds",
			requestBody: `{
				"username":"testuser",
				"password":"password123"
			}`,
			setupMock: func(ctx context.Context, mockSvc *mocks.Service) {
				mockSvc.
					On(
						"LoginUser",
						ctx,
						authDTO.LoginRequest{
							Username: "testuser",
							Password: "password123",
						},
					).
					Return("mocked-jwt-token", nil).
					Once()
			},
			expected: expected{
				statusCode:   http.StatusOK,
				bodyContains: "Logged in successfully",
			},
		},
		{
			name:        "should return 400 when request body is invalid JSON",
			requestBody: `{invalid json}`,
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
			name: `should return 400 when required field is missing`,
			requestBody: `{
				"username":"testuser"
			}`,
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
			name: "should return 401 when credentials are invalid",
			requestBody: `{
				"username":"testuser",
				"password":"wrong-password"
			}`,
			setupMock: func(
				ctx context.Context,
				mockSvc *mocks.Service,
			) {
				mockSvc.
					On(
						"LoginUser",
						ctx,
						authDTO.LoginRequest{
							Username: "testuser",
							Password: "wrong-password",
						},
					).
					Return(
						"",
						auth.ErrInvalidCredentials,
					).
					Once()
			},
			expected: expected{
				statusCode:   http.StatusUnauthorized,
				bodyContains: "Invalid username or password",
			},
		},
		{
			name: "should return 500 when service returns unexpected error",
			requestBody: `{
				"username":"testuser",
				"password":"password123"
			}`,
			setupMock: func(
				ctx context.Context,
				mockSvc *mocks.Service,
			) {
				mockSvc.
					On(
						"LoginUser",
						ctx,
						authDTO.LoginRequest{
							Username: "testuser",
							Password: "password123",
						},
					).
					Return(
						"",
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

			httpRequest := httptest.NewRequest(
				http.MethodPost,
				"/v1/users/login",
				strings.NewReader(tc.requestBody),
			)

			httpRequest.Header.Set(
				"Content-Type",
				"application/json",
			)

			ctx.Request = httpRequest

			tc.setupMock(ctx, mockSvc)

			handler := NewHandler(mockSvc)

			handler.Login(ctx)

			assert.Equal(
				t,
				tc.expected.statusCode,
				recorder.Code,
			)

			assert.Equal(
				t,
				"application/json; charset=utf-8",
				recorder.Header().Get("Content-Type"),
			)

			assert.Contains(
				t,
				recorder.Body.String(),
				tc.expected.bodyContains,
			)

			mockSvc.AssertExpectations(t)
		})
	}
}
