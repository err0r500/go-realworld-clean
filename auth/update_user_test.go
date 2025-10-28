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

func TestUpdateUser(t *testing.T) {
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
		Method: "PUT",
		Path:   "/api/user",
	}

	// Register a user to get a valid token
	registerReq := auth.RegisterReqBody{
		User: auth.RegisterUser{
			Username: "updateuser",
			Email:    "update@example.com",
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

	// Helper function to send authenticated update requests
	sendAuthenticatedUpdate := func(t *testing.T, updateBody auth.UpdateUserReqBody) (*http.Response, []byte) {
		t.Helper()
		bodyJSON, err := json.Marshal(updateBody)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", server.URL+"/api/user", strings.NewReader(string(bodyJSON)))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Token "+userToken)

		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		responseBody := readBody(t, response.Body)
		response.Body.Close()

		return response, responseBody
	}

	t.Run("update single fields", func(t *testing.T) {
		t.Run("update email", func(t *testing.T) {
			// given
			newEmail := "newemail@example.com"
			updateBody := auth.UpdateUserReqBody{
				User: auth.UpdateUserFields{
					Email: &newEmail,
				},
			}

			// when
			response, responseBody := sendAuthenticatedUpdate(t, updateBody)

			// then
			assert.Equal(t, http.StatusOK, response.StatusCode)
			responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

			var userResp auth.UserResponse
			require.NoError(t, json.Unmarshal(responseBody, &userResp))

			assert.Equal(t, newEmail, userResp.User.Email)
			assert.Equal(t, registerReq.User.Username, userResp.User.Username)
			assert.Equal(t, userToken, userResp.User.Token, "Token should be returned")
		})

		t.Run("update username", func(t *testing.T) {
			// given
			newUsername := "newusername"
			updateBody := auth.UpdateUserReqBody{
				User: auth.UpdateUserFields{
					Username: &newUsername,
				},
			}

			// when
			response, responseBody := sendAuthenticatedUpdate(t, updateBody)

			// then
			assert.Equal(t, http.StatusOK, response.StatusCode)
			responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

			var userResp auth.UserResponse
			require.NoError(t, json.Unmarshal(responseBody, &userResp))

			assert.Equal(t, newUsername, userResp.User.Username)
			assert.Equal(t, userToken, userResp.User.Token, "Token should be returned")
		})

		t.Run("update bio", func(t *testing.T) {
			// given
			newBio := "This is my new bio"
			updateBody := auth.UpdateUserReqBody{
				User: auth.UpdateUserFields{
					Bio: &newBio,
				},
			}

			// when
			response, responseBody := sendAuthenticatedUpdate(t, updateBody)

			// then
			assert.Equal(t, http.StatusOK, response.StatusCode)
			responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

			var userResp auth.UserResponse
			require.NoError(t, json.Unmarshal(responseBody, &userResp))

			assert.Equal(t, newBio, userResp.User.Bio)
			assert.Equal(t, userToken, userResp.User.Token, "Token should be returned")
		})

		t.Run("update image", func(t *testing.T) {
			// given
			newImage := "https://example.com/avatar.png"
			updateBody := auth.UpdateUserReqBody{
				User: auth.UpdateUserFields{
					Image: &newImage,
				},
			}

			// when
			response, responseBody := sendAuthenticatedUpdate(t, updateBody)

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
		// given
		newEmail := "multifieldupdate@example.com"
		newUsername := "multiuser"
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

		// when
		response, responseBody := sendAuthenticatedUpdate(t, updateBody)

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
		bodyJSON, _ := json.Marshal(updateBody)

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
		bodyJSON, _ := json.Marshal(updateBody)

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
				updateBody := auth.UpdateUserReqBody{
					User: auth.UpdateUserFields{},
				}
				bodyJSON, _ := json.Marshal(updateBody)

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
		// given
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
		// given - no fields to update
		updateBody := auth.UpdateUserReqBody{
			User: auth.UpdateUserFields{},
		}

		// when
		response, responseBody := sendAuthenticatedUpdate(t, updateBody)

		// then - should succeed with no changes
		assert.Equal(t, http.StatusOK, response.StatusCode)

		var userResp auth.UserResponse
		require.NoError(t, json.Unmarshal(responseBody, &userResp))
		assert.Equal(t, userToken, userResp.User.Token, "Token should be returned")
	})

	t.Run("update password and verify", func(t *testing.T) {
		// Register a new user for this test to avoid interfering with other tests
		newPassword := uuid.NewString()
		testUser := auth.RegisterReqBody{
			User: auth.RegisterUser{
				Username: "passworduser",
				Email:    "password@example.com",
				Password: newPassword,
			},
		}

		var testUserToken string
		{ // register test user
			regBody, _ := json.Marshal(testUser)
			regResp, err := http.Post(server.URL+"/api/users", "application/json", strings.NewReader(string(regBody)))
			require.NoError(t, err)
			require.Equal(t, http.StatusCreated, regResp.StatusCode)

			var regUserResp auth.UserResponse
			require.NoError(t, json.NewDecoder(regResp.Body).Decode(&regUserResp))
			testUserToken = regUserResp.User.Token
			regResp.Body.Close()
		}

		// Update password
		updatedPassword := uuid.NewString()
		updateBody := auth.UpdateUserReqBody{
			User: auth.UpdateUserFields{
				Password: &updatedPassword,
			},
		}

		bodyJSON, err := json.Marshal(updateBody)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", server.URL+"/api/user", strings.NewReader(string(bodyJSON)))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Token "+testUserToken)

		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)
		response.Body.Close()

		// Verify old password no longer works
		t.Run("old password should not work", func(t *testing.T) {
			loginBody := auth.LoginReqBody{
				User: auth.LoginUser{
					Email:    testUser.User.Email,
					Password: newPassword, // old password
				},
			}
			bodyJSON, _ := json.Marshal(loginBody)

			resp, err := http.Post(server.URL+"/api/users/login", "application/json", strings.NewReader(string(bodyJSON)))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		// Verify new password works
		t.Run("new password should work", func(t *testing.T) {
			loginBody := auth.LoginReqBody{
				User: auth.LoginUser{
					Email:    testUser.User.Email,
					Password: updatedPassword, // new password
				},
			}
			bodyJSON, _ := json.Marshal(loginBody)

			resp, err := http.Post(server.URL+"/api/users/login", "application/json", strings.NewReader(string(bodyJSON)))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var loginResp auth.UserResponse
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&loginResp))
			assert.NotEmpty(t, loginResp.User.Token)
		})
	})
}
