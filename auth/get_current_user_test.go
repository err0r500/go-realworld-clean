package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/err0r500/go-realworld-clean/driving/api/testutil"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCurrentUser(t *testing.T) {
	ctx := newAuthTestContext(t)
	defer ctx.cleanup()

	apiReq := testutil.APIRequest{
		Method: "GET",
		Path:   "/api/user",
	}

	t.Run("valid request with token returns user", func(t *testing.T) {
		// given
		username := "currentuser"
		email := "current@example.com"
		userToken, _ := ctx.givenUserExists(username, email)
		testutil.AuthenticatedRequestRespectsSpecs(t, t.Context(), ctx.specRouter, apiReq, userToken)

		// when
		response := ctx.whenGettingCurrentUser(userToken)
		defer response.Body.Close()

		// then
		userResp := assertSuccessfulGetCurrentUser(t, response)
		assert.Equal(t, username, userResp.User.Username)
		assert.Equal(t, email, userResp.User.Email)
		assert.Equal(t, "", userResp.User.Bio)
		assert.Equal(t, "", userResp.User.Image)
		assert.Equal(t, userToken, userResp.User.Token)

		ctx.ensureResponseMatchesSpec(t, apiReq, response)
	})

	t.Run("missing authorization header returns 401", func(t *testing.T) {
		// when
		response := ctx.whenGettingCurrentUserWithoutAuth()
		defer response.Body.Close()

		// then
		assertUnauthorized(t, response)
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		// when
		response := ctx.whenGettingCurrentUserWithHeader("Token invalid-token-here")
		defer response.Body.Close()

		// then
		assertUnauthorized(t, response)
	})

	t.Run("malformed authorization header returns 401", func(t *testing.T) {
		userToken, _ := ctx.givenUserExists("headeruser", "header@example.com")

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
				// when
				response := ctx.whenGettingCurrentUserWithHeader(tc.header)
				defer response.Body.Close()

				// then
				assertUnauthorized(t, response)
			})
		}
	})

	t.Run("openapi validation rejects malformed auth header", func(t *testing.T) {
		// given - request with invalid Authorization header format
		httpReq, err := http.NewRequest(apiReq.Method, apiReq.Path, nil)
		require.NoError(t, err)
		httpReq.Header.Set("Authorization", "Bearer invalid-format")

		route, pathParams, err := ctx.specRouter.FindRoute(httpReq)
		require.NoError(t, err)

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

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    httpReq,
			PathParams: pathParams,
			Route:      route,
			Options:    options,
		}

		// when
		err = openapi3filter.ValidateRequest(t.Context(), requestValidationInput)

		// then
		assert.Error(t, err, "OpenAPI validation should reject malformed Authorization header")
		assert.Contains(t, err.Error(), "invalid Authorization header format")
	})
}
