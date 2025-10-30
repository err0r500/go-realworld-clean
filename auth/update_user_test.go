package auth_test

import (
	"testing"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/err0r500/go-realworld-clean/driving/api/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	ctx := newAuthTestContext(t)
	defer ctx.cleanup()

	apiReq := testutil.APIRequest{
		Method: "PUT",
		Path:   "/api/user",
	}

	t.Run("update single fields", func(t *testing.T) {
		t.Run("update email", func(t *testing.T) {
			// given
			username := "updateemailuser"
			email := "updateemailuser@example.com"
			userToken, _ := ctx.givenUserExists(username, email)
			newEmail := "newemail@example.com"

			// when
			response := ctx.whenUpdatingUser(userToken).withEmail(newEmail).send()
			defer response.Body.Close()

			// then
			userResp := assertSuccessfulUpdate(t, response)
			assert.Equal(t, newEmail, userResp.User.Email)
			assert.Equal(t, username, userResp.User.Username)
			assert.Equal(t, userToken, userResp.User.Token)

			ctx.ensureResponseMatchesSpec(t, apiReq, response)
		})

		t.Run("update username", func(t *testing.T) {
			// given
			userToken, _ := ctx.givenUserExists("oldusername", "usernameuser@example.com")
			newUsername := "newusername"

			// when
			response := ctx.whenUpdatingUser(userToken).withUsername(newUsername).send()
			defer response.Body.Close()

			// then
			userResp := assertSuccessfulUpdate(t, response)
			assert.Equal(t, newUsername, userResp.User.Username)
			assert.Equal(t, userToken, userResp.User.Token)

			ctx.ensureResponseMatchesSpec(t, apiReq, response)
		})

		t.Run("update bio", func(t *testing.T) {
			// given
			userToken, _ := ctx.givenUserExists("biouser", "biouser@example.com")
			newBio := "This is my new bio"

			// when
			response := ctx.whenUpdatingUser(userToken).withBio(newBio).send()
			defer response.Body.Close()

			// then
			userResp := assertSuccessfulUpdate(t, response)
			assert.Equal(t, newBio, userResp.User.Bio)
			assert.Equal(t, userToken, userResp.User.Token)

			ctx.ensureResponseMatchesSpec(t, apiReq, response)
		})

		t.Run("update image", func(t *testing.T) {
			// given
			userToken, _ := ctx.givenUserExists("imageuser", "imageuser@example.com")
			newImage := "https://example.com/avatar.png"

			// when
			response := ctx.whenUpdatingUser(userToken).withImage(newImage).send()
			defer response.Body.Close()

			// then
			userResp := assertSuccessfulUpdate(t, response)
			assert.Equal(t, newImage, userResp.User.Image)
			assert.Equal(t, userToken, userResp.User.Token)

			ctx.ensureResponseMatchesSpec(t, apiReq, response)
		})
	})

	t.Run("update multiple fields", func(t *testing.T) {
		// given
		userToken, _ := ctx.givenUserExists("multiuser", "multiuser@example.com")
		newEmail := "multifieldupdate@example.com"
		newUsername := "updatedmultiuser"
		newBio := "Updated bio with multiple fields"
		newImage := "https://example.com/new-avatar.png"

		// when
		response := ctx.whenUpdatingUser(userToken).
			withEmail(newEmail).
			withUsername(newUsername).
			withBio(newBio).
			withImage(newImage).
			send()
		defer response.Body.Close()

		// then
		userResp := assertSuccessfulUpdate(t, response)
		assert.Equal(t, newEmail, userResp.User.Email)
		assert.Equal(t, newUsername, userResp.User.Username)
		assert.Equal(t, newBio, userResp.User.Bio)
		assert.Equal(t, newImage, userResp.User.Image)
		assert.Equal(t, userToken, userResp.User.Token)

		ctx.ensureResponseMatchesSpec(t, apiReq, response)
	})

	t.Run("missing authorization header returns 401", func(t *testing.T) {
		// when
		response := ctx.whenUpdatingUserWithoutAuth(auth.UpdateUserFields{})
		defer response.Body.Close()

		// then
		assertUnauthorized(t, response)
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		// when
		response := ctx.whenUpdatingUserWithHeader("Token invalid-token-here", auth.UpdateUserFields{})
		defer response.Body.Close()

		// then
		assertUnauthorized(t, response)
	})

	t.Run("malformed authorization header returns 401", func(t *testing.T) {
		testToken, _ := ctx.givenUserExists("headeruser", "headeruser@example.com")

		testCases := []struct {
			name   string
			header string
		}{
			{
				name:   "missing Token prefix",
				header: testToken,
			},
			{
				name:   "wrong prefix",
				header: "Bearer " + testToken,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// when
				response := ctx.whenUpdatingUserWithHeader(tc.header, auth.UpdateUserFields{})
				defer response.Body.Close()

				// then
				assertUnauthorized(t, response)
			})
		}
	})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		// given
		userToken, _ := ctx.givenUserExists("jsonuser", "jsonuser@example.com")

		// when
		response := ctx.whenUpdatingUserWithJSON(userToken, `{"user": {"email": "test@example.com"`)
		defer response.Body.Close()

		// then
		assertBadRequest(t, response)
	})

	t.Run("empty update request succeeds", func(t *testing.T) {
		// given
		userToken, _ := ctx.givenUserExists("emptyuser", "emptyuser@example.com")

		// when
		response := ctx.whenUpdatingUser(userToken).send()
		defer response.Body.Close()

		// then
		userResp := assertSuccessfulUpdate(t, response)
		assert.Equal(t, userToken, userResp.User.Token)
	})

	t.Run("update password allows login with new password", func(t *testing.T) {
		// given
		userToken, _ := ctx.givenUserExists("passworduser", "password@example.com")
		newPassword := uuid.NewString()

		// when
		response := ctx.whenUpdatingUser(userToken).withPassword(newPassword).send()
		defer response.Body.Close()

		// then
		userResp := assertSuccessfulUpdate(t, response)
		assert.Equal(t, userToken, userResp.User.Token)

		// verify can login with new password
		loginResp := ctx.whenLoggingInWith("password@example.com", newPassword)
		defer loginResp.Body.Close()
		assertSuccessfulLogin(t, loginResp)
	})

	t.Run("old password fails after password update", func(t *testing.T) {
		// given
		initialPassword := uuid.NewString()
		testEmail := "oldpasswordtest@example.com"

		regResp := ctx.whenRegisteringWith("oldpassworduser", testEmail, initialPassword)
		defer regResp.Body.Close()
		regUserResp := assertSuccessfulRegistration(t, regResp)
		userToken := regUserResp.User.Token

		newPassword := uuid.NewString()

		// when - update to new password
		updateResp := ctx.whenUpdatingUser(userToken).withPassword(newPassword).send()
		defer updateResp.Body.Close()
		assertSuccessfulUpdate(t, updateResp)

		// then - old password should not work
		loginResp := ctx.whenLoggingInWith(testEmail, initialPassword)
		defer loginResp.Body.Close()
		assertUnauthorized(t, loginResp)
	})
}
