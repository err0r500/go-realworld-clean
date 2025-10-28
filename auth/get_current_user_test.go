package auth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCurrentUser(t *testing.T) {
	kratosContainer := startKratosContainer(t.Context())
	defer func() {
		if err := kratosContainer.Terminate(t.Context()); err != nil {
			log.Fatal("failed to stop kratos container", err)
		}
	}()

	server := setupTestServer(t, getKratosURL(t.Context(), kratosContainer))
	defer server.Close()

	// Load OpenAPI spec for validation
	_, router := loadOpenAPISpec(t, t.Context())
	apiReq := APIRequest{
		Method: "GET",
		Path:   "/api/user",
	}

	// Register a user to get a valid token
	registerReq := auth.RegisterReqBody{
		User: auth.RegisterUser{
			Username: "currentuser",
			Email:    "current@example.com",
			Password: uuid.NewString(),
		},
	}

	var userToken string
	{ // register the user and capture token
		regBody, _ := json.Marshal(registerReq)
		regResp, err := http.Post(server.URL+"/api/users", "application/json", strings.NewReader(string(regBody)))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, regResp.StatusCode)

		var regUserResp auth.UserResponse
		require.NoError(t, json.NewDecoder(regResp.Body).Decode(&regUserResp))
		userToken = regUserResp.User.Token
		regResp.Body.Close()
	}

	t.Run("valid request with token", func(t *testing.T) {
		// given - validate request structure against OpenAPI spec
		authenticatedRequestRespectsSpecs(t, t.Context(), router, apiReq, userToken)

		// when - send the actual request
		req, err := http.NewRequest(apiReq.Method, server.URL+apiReq.Path, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Token "+userToken)

		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusOK, response.StatusCode)

		responseBody := readBody(t, response.Body)
		responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

		// Parse response and assert on values
		var userResp auth.UserResponse
		require.NoError(t, json.Unmarshal(responseBody, &userResp))

		assert.Equal(t, registerReq.User.Username, userResp.User.Username)
		assert.Equal(t, registerReq.User.Email, userResp.User.Email)
		assert.Equal(t, "", userResp.User.Bio)
		assert.Equal(t, "", userResp.User.Image)
		assert.Equal(t, userToken, userResp.User.Token, "Token should be returned in the response")
	})

	t.Run("missing authorization header", func(t *testing.T) {
		// given - request without Authorization header

		// when
		req, err := http.NewRequest(apiReq.Method, server.URL+apiReq.Path, nil)
		require.NoError(t, err)

		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("invalid token", func(t *testing.T) {
		// given - request with invalid token

		// when
		req, err := http.NewRequest(apiReq.Method, server.URL+apiReq.Path, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Token invalid-token-here")

		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("malformed authorization header", func(t *testing.T) {
		testCases := []struct {
			name   string
			header string
		}{
			{
				name:   "missing Token prefix",
				header: userToken,
			},
			{
				name:   "wrong prefix",
				header: "Bearer " + userToken,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// given
				req, err := http.NewRequest(apiReq.Method, server.URL+apiReq.Path, nil)
				require.NoError(t, err)
				req.Header.Set("Authorization", tc.header)

				// when
				response, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				defer response.Body.Close()

				// then
				assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
			})
		}
	})

	t.Run("openapi validation - rejects malformed auth header", func(t *testing.T) {
		// This test demonstrates that OpenAPI validation catches malformed headers
		// before they even reach the server

		// Create request with invalid Authorization header format (missing "Token" prefix)
		httpReq, err := http.NewRequest(apiReq.Method, apiReq.Path, nil)
		require.NoError(t, err)
		httpReq.Header.Set("Authorization", "Bearer invalid-format")

		// Find matching route
		route, pathParams, err := router.FindRoute(httpReq)
		require.NoError(t, err)

		// Configure validation options
		options := &openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
				if input.SecuritySchemeName == "Token" {
					authHeader := input.RequestValidationInput.Request.Header.Get("Authorization")
					if !strings.HasPrefix(authHeader, "Token ") {
						return fmt.Errorf("invalid Authorization header format, expected 'Token <jwt>'")
					}
					return nil
				}
				return fmt.Errorf("unknown security scheme: %s", input.SecuritySchemeName)
			},
		}

		// Validate request - should fail
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    httpReq,
			PathParams: pathParams,
			Route:      route,
			Options:    options,
		}

		err = openapi3filter.ValidateRequest(t.Context(), requestValidationInput)
		assert.Error(t, err, "OpenAPI validation should reject malformed Authorization header")
		assert.Contains(t, err.Error(), "invalid Authorization header format")
	})
}
