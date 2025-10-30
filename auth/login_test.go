package auth_test

import (
	"testing"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/err0r500/go-realworld-clean/driving/api/testutil"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	ctx := newAuthTestContext(t)
	defer ctx.cleanup()

	apiReq := testutil.APIRequest{
		Method: "POST",
		Path:   "/api/users/login",
	}

	t.Run("valid login returns user with token", func(t *testing.T) {
		// given
		username := "loginguser"
		email := "login@example.com"
		_, password := ctx.givenUserExists(username, email)

		apiReq.Body = toJSON(t, auth.LoginReqBody{
			User: auth.LoginUser{
				Email:    email,
				Password: password,
			},
		})
		ctx.ensureRequestMatchesSpec(t, apiReq)

		// when
		response := ctx.whenLoggingInWith(email, password)
		defer response.Body.Close()

		// then
		userResp := assertSuccessfulLogin(t, response)
		assert.Equal(t, email, userResp.User.Email)
		assert.Equal(t, username, userResp.User.Username)

		ctx.ensureResponseMatchesSpec(t, apiReq, response)
	})

	t.Run("invalid credentials return 401", func(t *testing.T) {
		_, password := ctx.givenUserExists("authuser", "auth@example.com")

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
				email:    "auth@example.com",
				password: "wrongpassword",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// when
				response := ctx.whenLoggingInWith(tc.email, tc.password)
				defer response.Body.Close()

				// then
				assertUnauthorized(t, response)
			})
		}
	})

	t.Run("missing required fields return 422", func(t *testing.T) {
		testCases := []struct {
			name string
			body string
		}{
			{
				name: "missing email",
				body: `{"user": {"password": "password123"}}`,
			},
			{
				name: "missing password",
				body: `{"user": {"email": "test@example.com"}}`,
			},
			{
				name: "empty user object",
				body: `{"user": {}}`,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// when
				response := ctx.whenLoggingInWithJSON(tc.body)
				defer response.Body.Close()

				// then
				assertUnprocessableEntity(t, response)
				apiReq.Body = tc.body
				ctx.ensureResponseMatchesSpec(t, apiReq, response)
			})
		}
	})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		// when
		response := ctx.whenLoggingInWithJSON(`{"user": {"email": "test@example.com"`)
		defer response.Body.Close()

		// then
		assertBadRequest(t, response)
	})
}
