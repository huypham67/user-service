package integration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserInfoEndpoint(t *testing.T) {
	t.Parallel()

	type expected struct {
		statusCode   int
		bodyContains string
	}

	testCases := []struct {
		name        string
		requestBody string
		setupAuth   func(t *testing.T, app *AuthenticatedTestApp, req *http.Request)
		expected    expected
	}{
		{
			name: "should return 200 when update user successfully",
			requestBody: `{
				"display_name": "Updated User",
				"email": "updated@example.com"
			}`,
			setupAuth: func(
				t *testing.T,
				app *AuthenticatedTestApp,
				req *http.Request,
			) {
				token, err := app.TokenGenerator.GenerateToken(
					"user-uuid-1",
					"testuser1",
					"testuser1@gmail.com",
				)

				require.NoError(t, err)

				req.Header.Set(
					"Authorization",
					"Bearer "+token,
				)
			},
			expected: expected{
				statusCode:   http.StatusOK,
				bodyContains: "Edit current user successfully!",
			},
		},
		{
			name: "should return 409 when email already exists",
			requestBody: `{
				"display_name": "Updated User",
				"email": "testuser2@gmail.com"
			}`,
			setupAuth: func(
				t *testing.T,
				app *AuthenticatedTestApp,
				req *http.Request,
			) {
				token, err := app.TokenGenerator.GenerateToken(
					"user-uuid-1",
					"testuser1",
					"testuser1@gmail.com",
				)

				require.NoError(t, err)

				req.Header.Set(
					"Authorization",
					"Bearer "+token,
				)
			},
			expected: expected{
				statusCode:   http.StatusConflict,
				bodyContains: "Email already exists",
			},
		},
		{
			name: "should return 500 when user does not exist",
			requestBody: `{
				"display_name": "Updated User",
				"email": "updated@example.com"
			}`,
			setupAuth: func(
				t *testing.T,
				app *AuthenticatedTestApp,
				req *http.Request,
			) {
				token, err := app.TokenGenerator.GenerateToken(
					"missing-user-id",
					"missing-user",
					"missinguser@gmail.com",
				)

				require.NoError(t, err)

				req.Header.Set(
					"Authorization",
					"Bearer "+token,
				)
			},
			expected: expected{
				statusCode:   http.StatusInternalServerError,
				bodyContains: "Internal Server Error",
			},
		},
		{
			name:        "should return 400 when request body is invalid JSON",
			requestBody: `{invalid json}`,
			setupAuth: func(
				t *testing.T,
				app *AuthenticatedTestApp,
				req *http.Request,
			) {
				token, err := app.TokenGenerator.GenerateToken(
					"user-uuid-1",
					"testuser1",
					"testuser1@gmail.com",
				)

				require.NoError(t, err)

				req.Header.Set(
					"Authorization",
					"Bearer "+token,
				)
			},
			expected: expected{
				statusCode:   http.StatusBadRequest,
				bodyContains: "Invalid request",
			},
		},
		{
			name: "should return 401 when authorization header is missing",
			requestBody: `{
				"display_name": "Updated User",
				"email": "updated@example.com"
			}`,
			setupAuth: func(
				t *testing.T,
				app *AuthenticatedTestApp,
				req *http.Request,
			) {
				// No auth header set
			},
			expected: expected{
				statusCode:   http.StatusUnauthorized,
				bodyContains: "missing authorization header",
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := setupProfileTestApp(t)

			httpRequest := httptest.NewRequest(
				http.MethodPut,
				"/api/user_service/v1/self/info",
				bytes.NewBufferString(tc.requestBody),
			)

			httpRequest.Header.Set(
				"Content-Type",
				"application/json",
			)

			tc.setupAuth(
				t,
				app,
				httpRequest,
			)

			httpRecorder := httptest.NewRecorder()

			app.Router.ServeHTTP(
				httpRecorder,
				httpRequest,
			)

			assert.Equal(
				t,
				tc.expected.statusCode,
				httpRecorder.Code,
			)

			assert.Contains(
				t,
				httpRecorder.Body.String(),
				tc.expected.bodyContains,
			)
		})
	}
}
