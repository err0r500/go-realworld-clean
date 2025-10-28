package auth_test

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	kratosContainer := startKratosContainer(t.Context())
	defer func() {
		if err := kratosContainer.Terminate(t.Context()); err != nil {
			log.Fatal("failed to stop kratos container", err)
		}
	}()

	server := setupTestServer(t, kratosContainer)
	defer server.Close()

	// Load OpenAPI spec for validation
	_, router := loadOpenAPISpec(t, t.Context())
	apiReq := APIRequest{
		Method: "POST",
		Path:   "/api/users",
	}

	t.Run("valid registration request", func(t *testing.T) {
		// given - create request using typed struct
		reqBody := auth.RegisterReqBody{
			User: auth.RegisterUser{
				Username: "testuser",
				Email:    "test@example.com",
				Password: uuid.NewString(),
			},
		}

		// Marshal to JSON for APIRequest
		bodyJSON, err := json.Marshal(reqBody)
		require.NoError(t, err)
		apiReq.Body = string(bodyJSON)

		requestRespectsSpecs(t, t.Context(), router, apiReq)

		// when
		response := sendAPIRequest(t, server.URL, apiReq)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusCreated, response.StatusCode)

		responseBody := readBody(t, response.Body)
		responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

		// Parse response and assert on values
		var userResp auth.UserResponse
		require.NoError(t, json.Unmarshal(responseBody, &userResp))

		assert.Equal(t, reqBody.User.Username, userResp.User.Username)
		assert.Equal(t, reqBody.User.Email, userResp.User.Email)
		assert.NotEmpty(t, userResp.User.Token, "should return a JWT token")
		assert.Equal(t, "", userResp.User.Bio)
		assert.Equal(t, "", userResp.User.Image)
	})

	t.Run("invalid request - missing required fields", func(t *testing.T) {
		testCases := []struct {
			name        string
			requestBody string
		}{
			{
				name:        "missing username",
				requestBody: `{"user": {"email": "test@example.com", "password": "password123"}}`,
			},
			{
				name:        "missing email",
				requestBody: `{"user": {"username": "testuser", "password": "password123"}}`,
			},
			{
				name:        "missing password",
				requestBody: `{"user": {"username": "testuser", "email": "test@example.com"}}`,
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
				responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, readBody(t, response.Body))
			})
		}
	})

	t.Run("invalid request - malformed JSON", func(t *testing.T) {
		// given
		apiReq.Body = `{"user": {"username": "testuser", "email": "test@example.com"`

		// when
		response := sendAPIRequest(t, server.URL, apiReq)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}
