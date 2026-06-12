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
	"github.com/huypham67/user-service/internal/model"
	"github.com/huypham67/user-service/internal/service/auth"
	"github.com/huypham67/user-service/internal/service/auth/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Register(t *testing.T) {
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
			name: "should return 201 when register succeeds",
			requestBody: `{
				"display_name":"Test User",
				"username":"testuser",
				"email":"test@example.com",
				"password":"password123"
			}`,
			setupMock: func(ctx context.Context, mockSvc *mocks.Service) {
				mockSvc.
					On(
						"RegisterUser",
						ctx,
						authDTO.RegisterUserRequest{
							DisplayName: "Test User",
							Username:    "testuser",
							Email:       "test@example.com",
							Password:    "password123",
						},
					).
					Return(
						&model.User{
							BaseModel: model.BaseModel{
								ID: "user-id",
							},
							DisplayName: "Test User",
							Username:    "testuser",
							Email:       "test@example.com",
						},
						nil,
					).
					Once()
			},
			expected: expected{
				statusCode:   http.StatusCreated,
				bodyContains: "Register an user successfully!",
			},
		},
		{
			name:        "should return 400 when request body is invalid JSON",
			requestBody: `{invalid json}`,
			setupMock:   func(ctx context.Context, mockSvc *mocks.Service) {},
			expected: expected{
				statusCode:   http.StatusBadRequest,
				bodyContains: "Invalid request body",
			},
		},
		{
			name: "should return 400 when required field is missing",
			requestBody: `{
				"display_name":"Test User",
				"username":"testuser"
			}`,
			setupMock: func(ctx context.Context, mockSvc *mocks.Service) {},
			expected: expected{
				statusCode:   http.StatusBadRequest,
				bodyContains: "Invalid request body",
			},
		},
		{
			name: "should return 409 when user already exists",
			requestBody: `{
				"display_name":"Test User",
				"username":"existinguser",
				"email":"existing@example.com",
				"password":"password123"
			}`,
			setupMock: func(ctx context.Context, mockSvc *mocks.Service) {
				mockSvc.
					On(
						"RegisterUser",
						ctx,
						authDTO.RegisterUserRequest{
							DisplayName: "Test User",
							Username:    "existinguser",
							Email:       "existing@example.com",
							Password:    "password123",
						},
					).
					Return(nil, auth.ErrUserAlreadyExists).
					Once()
			},
			expected: expected{
				statusCode:   http.StatusConflict,
				bodyContains: "User already exists",
			},
		},
		{
			name: "should return 500 when service returns unexpected error",
			requestBody: `{
				"display_name":"Test User",
				"username":"testuser",
				"email":"test@example.com",
				"password":"password123"
			}`,
			setupMock: func(ctx context.Context, mockSvc *mocks.Service) {
				mockSvc.
					On(
						"RegisterUser",
						ctx,
						authDTO.RegisterUserRequest{
							DisplayName: "Test User",
							Username:    "testuser",
							Email:       "test@example.com",
							Password:    "password123",
						},
					).
					Return(nil, errors.New("database error")).
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
				"/v1/users/register",
				strings.NewReader(tc.requestBody),
			)

			httpRequest.Header.Set(
				"Content-Type",
				"application/json",
			)

			ctx.Request = httpRequest

			tc.setupMock(ctx, mockSvc)

			handler := NewHandler(mockSvc)

			handler.Register(ctx)

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
