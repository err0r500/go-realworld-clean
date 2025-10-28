package auth_test

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
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
		Path:   "/api/users/login",
	}

	registerReq := auth.RegisterReqBody{
		User: auth.RegisterUser{
			Username: "loginuser",
			Email:    "login@example.com",
			Password: uuid.NewString(),
		},
	}

	{ // register the user
		regBody, _ := json.Marshal(registerReq)
		regResp, err := http.Post(server.URL+"/api/users", "application/json", strings.NewReader(string(regBody)))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, regResp.StatusCode)
		regResp.Body.Close()
	}

	t.Run("valid login request", func(t *testing.T) {
		// given - login with registered credentials
		loginBody := auth.LoginReqBody{
			User: auth.LoginUser{
				Email:    registerReq.User.Email,
				Password: registerReq.User.Password,
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
		assert.Equal(t, registerReq.User.Username, userResp.User.Username)
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
				password: registerReq.User.Password,
			},
			{
				name:     "wrong password",
				email:    registerReq.User.Email,
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
				responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, readBody(t, response.Body))
			})
		}
	})

	t.Run("invalid request - malformed JSON", func(t *testing.T) {
		// given
		apiReq.Body = `{"user": {"email": "test@example.com"`

		// when
		response := sendAPIRequest(t, server.URL, apiReq)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}
