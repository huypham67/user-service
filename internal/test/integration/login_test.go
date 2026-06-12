package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	authDTO "github.com/huypham67/user-service/internal/dto/auth"
	"github.com/huypham67/user-service/internal/test/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginEndpoint(t *testing.T) {
	t.Parallel()

	type expected struct {
		statusCode   int
		bodyContains string
	}

	testCases := []struct {
		name        string
		requestBody string
		expected    expected
	}{
		{
			name: "should return 200 when login is successful",
			requestBody: fmt.Sprintf(
				`{
					"username": "%s",
					"password": "%s"
				}`,
				"testuser1",
				fixtures.TestPassword,
			),
			expected: expected{
				statusCode:   http.StatusOK,
				bodyContains: "Logged in successfully!",
			},
		},
		{
			name: "should return 401 when username does not exist",
			requestBody: `{
				"username": "notfound",
				"password": "password123"
			}`,
			expected: expected{
				statusCode:   http.StatusUnauthorized,
				bodyContains: "Invalid username or password",
			},
		},
		{
			name: "should return 401 when password is incorrect",
			requestBody: `{
				"username": "testuser1",
				"password": "wrongpassword"
			}`,
			expected: expected{
				statusCode:   http.StatusUnauthorized,
				bodyContains: "Invalid username or password",
			},
		},
		{
			name:        "should return 400 when request body is invalid JSON",
			requestBody: `{invalid json}`,
			expected: expected{
				statusCode:   http.StatusBadRequest,
				bodyContains: "Invalid request body",
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := setupAuthTestApp(t)

			httpRequest := httptest.NewRequest(
				http.MethodPost,
				"/api/user_service/v1/users/login",
				bytes.NewBufferString(tc.requestBody),
			)

			httpRequest.Header.Set(
				"Content-Type",
				"application/json",
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
				var resp authDTO.LoginResponse

				err := json.Unmarshal(
					httpRecorder.Body.Bytes(),
					&resp,
				)

				require.NoError(t, err)

				assert.NotEmpty(
					t,
					resp.Data,
				)

				assert.Equal(
					t,
					"Logged in successfully!",
					resp.Message,
				)
			}
		})
	}
}
