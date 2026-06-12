package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	profileDTO "github.com/huypham67/user-service/internal/dto/profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserInfoEndpoint(t *testing.T) {
	t.Parallel()

	type expected struct {
		statusCode   int
		bodyContains string
	}

	testCases := []struct {
		name      string
		setupAuth func(
			t *testing.T,
			app *AuthenticatedTestApp,
			req *http.Request,
		)
		expected expected
	}{
		{
			name: "should return 200 when user exists",
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
				bodyContains: "User information retrieved successfully!",
			},
		},
		{
			name: "should return 500 when user does not exist",
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
			name: "should return 401 when authorization header is missing",
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
				http.MethodGet,
				"/api/user_service/v1/self/info",
				nil,
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

			if tc.expected.statusCode == http.StatusOK {
				var resp profileDTO.UserResponse

				err := json.Unmarshal(
					httpRecorder.Body.Bytes(),
					&resp,
				)

				require.NoError(t, err)

				assert.Equal(
					t,
					"user-uuid-1",
					resp.Data.ID,
				)

				assert.Equal(
					t,
					"testuser1@gmail.com",
					resp.Data.Email,
				)

				assert.Equal(
					t,
					"User information retrieved successfully!",
					resp.Message,
				)
			}
		})
	}
}
