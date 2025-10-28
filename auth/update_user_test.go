package auth_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateUser(t *testing.T) {
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
		Method: "PUT",
		Path:   "/api/user",
	}

	// Note: Not using a shared user here - each subtest creates its own user for test isolation

	t.Run("update single fields", func(t *testing.T) {
		t.Run("update email", func(t *testing.T) {
			// given - create dedicated user for this test
			userToken, _ := registerTestUser(t, server.URL, "emailuser", "emailuser@example.com")

			newEmail := "newemail@example.com"
			updateBody := auth.UpdateUserReqBody{
				User: auth.UpdateUserFields{
					Email: &newEmail,
				},
			}

			bodyJSON, err := json.Marshal(updateBody)
			require.NoError(t, err)

			// when
			response, responseBody := sendAuthenticatedRequest(t, "PUT", server.URL+"/api/user", string(bodyJSON), userToken)

			// then
			assert.Equal(t, http.StatusOK, response.StatusCode)
			responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

			var userResp auth.UserResponse
			require.NoError(t, json.Unmarshal(responseBody, &userResp))

			assert.Equal(t, newEmail, userResp.User.Email)
			assert.Equal(t, "emailuser", userResp.User.Username)
			assert.Equal(t, userToken, userResp.User.Token, "Token should be returned")
		})

		t.Run("update username", func(t *testing.T) {
			// given - create dedicated user for this test
			userToken, _ := registerTestUser(t, server.URL, "oldusername", "usernameuser@example.com")

			newUsername := "newusername"
			updateBody := auth.UpdateUserReqBody{
				User: auth.UpdateUserFields{
					Username: &newUsername,
				},
			}

			bodyJSON, err := json.Marshal(updateBody)
			require.NoError(t, err)

			// when
			response, responseBody := sendAuthenticatedRequest(t, "PUT", server.URL+"/api/user", string(bodyJSON), userToken)

			// then
			assert.Equal(t, http.StatusOK, response.StatusCode)
			responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

			var userResp auth.UserResponse
			require.NoError(t, json.Unmarshal(responseBody, &userResp))

			assert.Equal(t, newUsername, userResp.User.Username)
			assert.Equal(t, userToken, userResp.User.Token, "Token should be returned")
		})

		t.Run("update bio", func(t *testing.T) {
			// given - create dedicated user for this test
			userToken, _ := registerTestUser(t, server.URL, "biouser", "biouser@example.com")

			newBio := "This is my new bio"
			updateBody := auth.UpdateUserReqBody{
				User: auth.UpdateUserFields{
					Bio: &newBio,
				},
			}

			bodyJSON, err := json.Marshal(updateBody)
			require.NoError(t, err)

			// when
			response, responseBody := sendAuthenticatedRequest(t, "PUT", server.URL+"/api/user", string(bodyJSON), userToken)

			// then
			assert.Equal(t, http.StatusOK, response.StatusCode)
			responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

			var userResp auth.UserResponse
			require.NoError(t, json.Unmarshal(responseBody, &userResp))

			assert.Equal(t, newBio, userResp.User.Bio)
			assert.Equal(t, userToken, userResp.User.Token, "Token should be returned")
		})

		t.Run("update image", func(t *testing.T) {
			// given - create dedicated user for this test
			userToken, _ := registerTestUser(t, server.URL, "imageuser", "imageuser@example.com")

			newImage := "https://example.com/avatar.png"
			updateBody := auth.UpdateUserReqBody{
				User: auth.UpdateUserFields{
					Image: &newImage,
				},
			}

			bodyJSON, err := json.Marshal(updateBody)
			require.NoError(t, err)

			// when
			response, responseBody := sendAuthenticatedRequest(t, "PUT", server.URL+"/api/user", string(bodyJSON), userToken)

			// then
			assert.Equal(t, http.StatusOK, response.StatusCode)
			responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

			var userResp auth.UserResponse
			require.NoError(t, json.Unmarshal(responseBody, &userResp))

			assert.Equal(t, newImage, userResp.User.Image)
			assert.Equal(t, userToken, userResp.User.Token, "Token should be returned")
		})
	})

	t.Run("update multiple fields", func(t *testing.T) {
		// given - create dedicated user for this test
		userToken, _ := registerTestUser(t, server.URL, "multiuser", "multiuser@example.com")

		newEmail := "multifieldupdate@example.com"
		newUsername := "updatedmultiuser"
		newBio := "Updated bio with multiple fields"
		newImage := "https://example.com/new-avatar.png"

		updateBody := auth.UpdateUserReqBody{
			User: auth.UpdateUserFields{
				Email:    &newEmail,
				Username: &newUsername,
				Bio:      &newBio,
				Image:    &newImage,
			},
		}

		bodyJSON, err := json.Marshal(updateBody)
		require.NoError(t, err)

		// when
		response, responseBody := sendAuthenticatedRequest(t, "PUT", server.URL+"/api/user", string(bodyJSON), userToken)

		// then
		assert.Equal(t, http.StatusOK, response.StatusCode)
		responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

		var userResp auth.UserResponse
		require.NoError(t, json.Unmarshal(responseBody, &userResp))

		assert.Equal(t, newEmail, userResp.User.Email)
		assert.Equal(t, newUsername, userResp.User.Username)
		assert.Equal(t, newBio, userResp.User.Bio)
		assert.Equal(t, newImage, userResp.User.Image)
		assert.Equal(t, userToken, userResp.User.Token, "Token should be returned")
	})

	t.Run("missing authorization header", func(t *testing.T) {
		// given
		updateBody := auth.UpdateUserReqBody{
			User: auth.UpdateUserFields{},
		}
		bodyJSON, err := json.Marshal(updateBody)
		require.NoError(t, err)

		// when
		req, err := http.NewRequest("PUT", server.URL+"/api/user", strings.NewReader(string(bodyJSON)))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("invalid token", func(t *testing.T) {
		// given
		updateBody := auth.UpdateUserReqBody{
			User: auth.UpdateUserFields{},
		}
		bodyJSON, err := json.Marshal(updateBody)
		require.NoError(t, err)

		// when
		req, err := http.NewRequest("PUT", server.URL+"/api/user", strings.NewReader(string(bodyJSON)))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Token invalid-token-here")

		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("malformed authorization header", func(t *testing.T) {
		// Create a user to get a real token for these tests
		testToken, _ := registerTestUser(t, server.URL, "headeruser", "headeruser@example.com")

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
				// given
				updateBody := auth.UpdateUserReqBody{
					User: auth.UpdateUserFields{},
				}
				bodyJSON, err := json.Marshal(updateBody)
				require.NoError(t, err)

				req, err := http.NewRequest("PUT", server.URL+"/api/user", strings.NewReader(string(bodyJSON)))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
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

	t.Run("invalid request - malformed JSON", func(t *testing.T) {
		// given - create dedicated user for this test
		userToken, _ := registerTestUser(t, server.URL, "jsonuser", "jsonuser@example.com")

		// Intentionally malformed JSON to test error handling
		req, err := http.NewRequest("PUT", server.URL+"/api/user", strings.NewReader(`{"user": {"email": "test@example.com"`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Token "+userToken)

		// when
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("empty update request", func(t *testing.T) {
		// given - create dedicated user for this test
		userToken, _ := registerTestUser(t, server.URL, "emptyuser", "emptyuser@example.com")

		// No fields to update
		updateBody := auth.UpdateUserReqBody{
			User: auth.UpdateUserFields{},
		}

		bodyJSON, err := json.Marshal(updateBody)
		require.NoError(t, err)

		// when
		response, responseBody := sendAuthenticatedRequest(t, "PUT", server.URL+"/api/user", string(bodyJSON), userToken)

		// then - should succeed with no changes
		assert.Equal(t, http.StatusOK, response.StatusCode)

		var userResp auth.UserResponse
		require.NoError(t, json.Unmarshal(responseBody, &userResp))
		assert.Equal(t, userToken, userResp.User.Token, "Token should be returned")
	})

	t.Run("update password", func(t *testing.T) {
		// given - create dedicated user for password update test
		userToken, _ := registerTestUser(t, server.URL, "passworduser", "password@example.com")

		// Update password
		updatedPassword := uuid.NewString()
		updateBody := auth.UpdateUserReqBody{
			User: auth.UpdateUserFields{
				Password: &updatedPassword,
			},
		}

		bodyJSON, err := json.Marshal(updateBody)
		require.NoError(t, err)

		// when
		response, responseBody := sendAuthenticatedRequest(t, "PUT", server.URL+"/api/user", string(bodyJSON), userToken)

		// then
		require.Equal(t, http.StatusOK, response.StatusCode)

		var userResp auth.UserResponse
		require.NoError(t, json.Unmarshal(responseBody, &userResp))
		assert.Equal(t, userToken, userResp.User.Token, "Token should be returned")

		// Verify old password no longer works (stored in registerTestUser, but we can't retrieve it)
		// Note: We can't test old password here because registerTestUser doesn't return it
		// This demonstrates a limitation - consider returning password from registerTestUser if needed

		// Verify new password works
		loginBody := auth.LoginReqBody{
			User: auth.LoginUser{
				Email:    "password@example.com",
				Password: updatedPassword,
			},
		}
		loginBodyJSON, err := json.Marshal(loginBody)
		require.NoError(t, err)

		loginResp, err := http.Post(server.URL+"/api/users/login", "application/json", strings.NewReader(string(loginBodyJSON)))
		require.NoError(t, err)
		defer loginResp.Body.Close()

		assert.Equal(t, http.StatusOK, loginResp.StatusCode, "Should be able to login with new password")

		var loginUserResp auth.UserResponse
		require.NoError(t, json.NewDecoder(loginResp.Body).Decode(&loginUserResp))
		assert.NotEmpty(t, loginUserResp.User.Token)
	})

	t.Run("old password fails after password update", func(t *testing.T) {
		// given - register user with known password
		initialPassword := uuid.NewString()
		testEmail := "oldpasswordtest@example.com"

		// Manually register to keep track of initial password
		registerReq := auth.RegisterReqBody{
			User: auth.RegisterUser{
				Username: "oldpassworduser",
				Email:    testEmail,
				Password: initialPassword,
			},
		}

		regBody, err := json.Marshal(registerReq)
		require.NoError(t, err)

		regResp, err := http.Post(server.URL+"/api/users", "application/json", strings.NewReader(string(regBody)))
		require.NoError(t, err)
		defer regResp.Body.Close()
		require.Equal(t, http.StatusCreated, regResp.StatusCode)

		var regUserResp auth.UserResponse
		require.NoError(t, json.NewDecoder(regResp.Body).Decode(&regUserResp))
		userToken := regUserResp.User.Token

		// Update to new password
		newPassword := uuid.NewString()
		updateBody := auth.UpdateUserReqBody{
			User: auth.UpdateUserFields{
				Password: &newPassword,
			},
		}

		updateBodyJSON, err := json.Marshal(updateBody)
		require.NoError(t, err)

		updateResp, _ := sendAuthenticatedRequest(t, "PUT", server.URL+"/api/user", string(updateBodyJSON), userToken)
		require.Equal(t, http.StatusOK, updateResp.StatusCode)

		// when - try to login with old password
		loginBody := auth.LoginReqBody{
			User: auth.LoginUser{
				Email:    testEmail,
				Password: initialPassword, // old password
			},
		}
		loginBodyJSON, err := json.Marshal(loginBody)
		require.NoError(t, err)

		loginResp, err := http.Post(server.URL+"/api/users/login", "application/json", strings.NewReader(string(loginBodyJSON)))
		require.NoError(t, err)
		defer loginResp.Body.Close()

		// then - should fail
		assert.Equal(t, http.StatusUnauthorized, loginResp.StatusCode, "Old password should not work after update")
	})
}
