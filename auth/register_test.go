package auth_test

import (
	"testing"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/err0r500/go-realworld-clean/driving/api/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	ctx := newAuthTestContext(t)
	defer ctx.cleanup()

	apiReq := testutil.APIRequest{
		Method: "POST",
		Path:   "/api/users",
	}

	t.Run("valid registration returns user with token", func(t *testing.T) {
		// given
		reqBody := auth.RegisterReqBody{
			User: auth.RegisterUser{
				Username: "testuser",
				Email:    "test@example.com",
				Password: uuid.NewString(),
			},
		}
		apiReq.Body = toJSON(t, reqBody)
		ctx.ensureRequestMatchesSpec(t, apiReq)

		// when
		response := ctx.whenRegisteringWith(reqBody.User.Username, reqBody.User.Email, reqBody.User.Password)
		defer response.Body.Close()

		// then
		userResp := assertSuccessfulRegistration(t, response)
		assert.Equal(t, reqBody.User.Username, userResp.User.Username)
		assert.Equal(t, reqBody.User.Email, userResp.User.Email)
		assert.Equal(t, "", userResp.User.Bio)
		assert.Equal(t, "", userResp.User.Image)

		ctx.ensureResponseMatchesSpec(t, apiReq, response)
	})

	t.Run("missing required fields return 422", func(t *testing.T) {
		testCases := []struct {
			name string
			body string
		}{
			{
				name: "missing username",
				body: `{"user": {"email": "test@example.com", "password": "password123"}}`,
			},
			{
				name: "missing email",
				body: `{"user": {"username": "testuser", "password": "password123"}}`,
			},
			{
				name: "missing password",
				body: `{"user": {"username": "testuser", "email": "test@example.com"}}`,
			},
			{
				name: "empty user object",
				body: `{"user": {}}`,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// when
				response := ctx.whenRegisteringWithJSON(tc.body)
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
		response := ctx.whenRegisteringWithJSON(`{"user": {"username": "testuser", "email": "test@example.com"`)
		defer response.Body.Close()

		// then
		assertBadRequest(t, response)
	})
}
