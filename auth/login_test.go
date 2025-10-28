package auth_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	kratosContainer, err := startKratosContainer(t.Context())
	require.NoError(t, err, "Failed to start Kratos container")
	defer func() {
		if err := kratosContainer.Terminate(t.Context()); err != nil {
			t.Errorf("failed to stop kratos container: %v", err)
		}
	}()

	server := setupTestServer(t, kratosContainer)
	defer server.Close()

	// Load OpenAPI spec for validation
	_, router := loadOpenAPISpec(t, t.Context())
	apiReq := APIRequest{
		Method: "POST",
		Path:   "/api/users/login",
	}

	// Register a test user for authentication tests
	_, password := registerTestUser(t, server.URL, "loginuser", "login@example.com")
	testEmail := "login@example.com"

	t.Run("valid login request", func(t *testing.T) {
		// given - login with registered credentials
		loginBody := auth.LoginReqBody{
			User: auth.LoginUser{
				Email:    testEmail,
				Password: password,
			},
		}

		bodyJSON, err := json.Marshal(loginBody)
		require.NoError(t, err)
		apiReq.Body = string(bodyJSON)

		requestRespectsSpecs(t, t.Context(), router, apiReq)

		// when
		response := sendAPIRequest(t, server.URL, apiReq)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusOK, response.StatusCode)

		responseBody := readBody(t, response.Body)
		responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

		// Parse response and assert on values
		var userResp auth.UserResponse
		require.NoError(t, json.Unmarshal(responseBody, &userResp))

		assert.Equal(t, loginBody.User.Email, userResp.User.Email)
		assert.NotEmpty(t, userResp.User.Token, "should return a JWT token")
		assert.Equal(t, "loginuser", userResp.User.Username)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		testCases := []struct {
			name     string
			email    string
			password string
		}{
			{
				name:     "non-existent user",
				email:    "nonexistent@example.com",
				password: password,
			},
			{
				name:     "wrong password",
				email:    testEmail,
				password: "wrongpassword",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// given
				loginBody := auth.LoginReqBody{
					User: auth.LoginUser{
						Email:    tc.email,
						Password: tc.password,
					},
				}
				bodyJSON, _ := json.Marshal(loginBody)
				apiReq.Body = string(bodyJSON)

				// when
				response := sendAPIRequest(t, server.URL, apiReq)
				defer response.Body.Close()

				// then
				assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
			})
		}
	})

	t.Run("invalid request - missing required fields", func(t *testing.T) {
		// Note: Using raw JSON strings here to test OpenAPI validation of incomplete payloads
		testCases := []struct {
			name        string
			requestBody string
		}{
			{
				name:        "missing email",
				requestBody: `{"user": {"password": "password123"}}`,
			},
			{
				name:        "missing password",
				requestBody: `{"user": {"email": "test@example.com"}}`,
			},
			{
				name:        "empty user object",
				requestBody: `{"user": {}}`,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// given
				apiReq.Body = tc.requestBody

				// when
				response := sendAPIRequest(t, server.URL, apiReq)
				defer response.Body.Close()

				// then
				assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
				responseBody := readBody(t, response.Body)
				responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)
			})
		}
	})

	t.Run("invalid request - malformed JSON", func(t *testing.T) {
		// given - intentionally malformed JSON to test error handling
		apiReq.Body = `{"user": {"email": "test@example.com"`

		// when
		response := sendAPIRequest(t, server.URL, apiReq)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}
