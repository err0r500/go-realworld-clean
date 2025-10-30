package profile_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/driving/api/profile"
	"github.com/err0r500/go-realworld-clean/driving/api/testutil"
	"github.com/err0r500/go-realworld-clean/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestGetProfileEndpoint(t *testing.T) {
	// Load OpenAPI spec for validation
	_, specsRouter := testutil.LoadOpenAPISpec(t, context.Background(), "../../../_specs/openapi.yml")

	t.Run("valid unauthenticated request returns profile", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			existingProfile := anyProfile(t)
			mux := http.NewServeMux()

			attachGetProfileFnToHttpHandler(mux, func(ctx context.Context, requestingUserID *string, targetUsername string) (*domain.Profile, error) {
				assert.Nil(t, requestingUserID)                           // unauthenticated requests don't have a requestingUser
				assert.Equal(t, existingProfile.Username, targetUsername) // in this scenario, we try to retrieve the existing profile
				return &existingProfile, nil
			})

			requestParams := newGetProfileReqParams(existingProfile.Username)

			// when
			response := getProfile(mux, requestParams, nil)

			// then
			assert.Equal(t, http.StatusOK, response.StatusCode)

			respBody := parseResponseBody(tt, testutil.CloneResponse(tt, response))
			assert.Equal(t, existingProfile.Username, respBody.Profile.Username)
			assert.Equal(t, existingProfile.Bio, respBody.Profile.Bio)
			assert.Equal(t, existingProfile.Image, respBody.Profile.Image)
			assert.Equal(t, existingProfile.Following, respBody.Profile.Following) // doesn't make much sense without being authenticated but it's not the responsibility of this layer

			testutil.EnsureRequestAndResponseMatchSpecs(tt, tt.Context(), specsRouter, requestParams, response)
		})
	})

	t.Run("valid authenticated request passes the requesting user to business logic", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			existingProfile := anyProfile(t)
			requestFromUserId := nonEmptyString().Draw(t, "requestingUserId")
			mux := http.NewServeMux()

			attachGetProfileFnToHttpHandler(mux, func(ctx context.Context, requestingUserID *string, targetUsername string) (*domain.Profile, error) {
				assert.Equal(t, requestFromUserId, *requestingUserID) // called with the requestingUserId
				assert.Equal(t, existingProfile.Username, targetUsername)
				return &existingProfile, nil
			})

			requestParams := newGetProfileReqParams(existingProfile.Username)

			// when
			response := getProfile(mux, requestParams, &requestFromUserId)

			// then
			assert.Equal(t, http.StatusOK, response.StatusCode)

			testutil.EnsureRequestAndResponseMatchSpecs(tt, tt.Context(), specsRouter, requestParams, response)
		})
	})

	t.Run("user not found returns 404", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			mux := http.NewServeMux()

			attachGetProfileFnToHttpHandler(mux, func(ctx context.Context, requestingUserID *string, targetUsername string) (*domain.Profile, error) {
				return nil, usecase.ErrUserNotFound
			})

			requestParams := newGetProfileReqParams(nonEmptyString().Draw(t, "nonExistentUser"))

			// when
			response := getProfile(mux, requestParams, nil)

			// then
			assert.Equal(t, http.StatusNotFound, response.StatusCode)

			testutil.EnsureRequestAndResponseMatchSpecs(tt, tt.Context(), specsRouter, requestParams, response)
		})
	})

	t.Run("internal error returns 422", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			mux := http.NewServeMux()

			attachGetProfileFnToHttpHandler(mux, func(ctx context.Context, requestingUserID *string, targetUsername string) (*domain.Profile, error) {
				return nil, errors.New(nonEmptyString().Draw(t, "any unknown error"))
			})

			requestParams := newGetProfileReqParams(nonEmptyString().Draw(t, "nonExistentUser"))

			// when
			response := getProfile(mux, requestParams, nil)

			// then
			assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
			testutil.EnsureRequestAndResponseMatchSpecs(tt, tt.Context(), specsRouter, requestParams, response)
		})
	})
}

func TestFollowUserEndpoint(t *testing.T) {
	// Load OpenAPI spec for validation
	_, specsRouter := testutil.LoadOpenAPISpec(t, context.Background(), "../../../_specs/openapi.yml")

	t.Run("valid authenticated request returns profile with following=true", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			targetProfile := anyProfile(t)
			targetProfile.Following = true // ensure following is true
			followerID := nonEmptyString().Draw(t, "followerID")
			mux := http.NewServeMux()

			attachFollowUserFnToHttpHandler(mux, func(ctx context.Context, followerIDArg, followeeUsername string) (*domain.Profile, error) {
				assert.Equal(t, followerID, followerIDArg)
				assert.Equal(t, targetProfile.Username, followeeUsername)
				return &targetProfile, nil
			})

			requestParams := newFollowUserReqParams(targetProfile.Username)

			// when
			response := followUser(mux, requestParams, &followerID)

			// then
			assert.Equal(t, http.StatusOK, response.StatusCode)

			respBody := parseResponseBody(tt, testutil.CloneResponse(tt, response))
			assert.Equal(t, targetProfile.Username, respBody.Profile.Username)
			assert.Equal(t, targetProfile.Bio, respBody.Profile.Bio)
			assert.Equal(t, targetProfile.Image, respBody.Profile.Image)
			assert.True(t, respBody.Profile.Following)

			testutil.EnsureRespMatchesSpecs(tt, tt.Context(), specsRouter, requestParams, response)
		})
	})

	t.Run("missing authentication returns 401", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			mux := http.NewServeMux()

			attachFollowUserFnToHttpHandler(mux, func(ctx context.Context, followerID, followeeUsername string) (*domain.Profile, error) {
				return nil, errors.New("should not be called")
			})

			requestParams := newFollowUserReqParams(nonEmptyString().Draw(t, "username"))

			// when
			response := followUser(mux, requestParams, nil)

			// then
			assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		})
	})

	t.Run("user not found returns 404", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			mux := http.NewServeMux()

			attachFollowUserFnToHttpHandler(mux, func(ctx context.Context, followerID, followeeUsername string) (*domain.Profile, error) {
				return nil, usecase.ErrUserNotFound
			})

			requestParams := newFollowUserReqParams(nonEmptyString().Draw(t, "nonExistentUser"))

			// when
			followerID := nonEmptyString().Draw(t, "followerID")
			response := followUser(mux, requestParams, &followerID)

			// then
			assert.Equal(t, http.StatusNotFound, response.StatusCode)

			testutil.EnsureRespMatchesSpecs(tt, tt.Context(), specsRouter, requestParams, response)
		})
	})

	t.Run("cannot follow yourself returns 422", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			mux := http.NewServeMux()

			attachFollowUserFnToHttpHandler(mux, func(ctx context.Context, followerID, followeeUsername string) (*domain.Profile, error) {
				return nil, usecase.ErrCannotFollowYourself
			})

			requestParams := newFollowUserReqParams(nonEmptyString().Draw(t, "username"))

			// when
			followerID := nonEmptyString().Draw(t, "followerID")
			response := followUser(mux, requestParams, &followerID)

			// then
			assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)

			testutil.EnsureRespMatchesSpecs(tt, tt.Context(), specsRouter, requestParams, response)
		})
	})

	t.Run("internal error returns 500", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			mux := http.NewServeMux()

			attachFollowUserFnToHttpHandler(mux, func(ctx context.Context, followerID, followeeUsername string) (*domain.Profile, error) {
				return nil, errors.New(nonEmptyString().Draw(t, "any unknown error"))
			})

			requestParams := newFollowUserReqParams(nonEmptyString().Draw(t, "username"))

			// when
			followerID := nonEmptyString().Draw(t, "followerID")
			response := followUser(mux, requestParams, &followerID)

			// then
			assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)

			testutil.EnsureRespMatchesSpecs(tt, tt.Context(), specsRouter, requestParams, response)
		})
	})
}

func TestUnfollowUserEndpoint(t *testing.T) {
	// Load OpenAPI spec for validation
	_, specsRouter := testutil.LoadOpenAPISpec(t, context.Background(), "../../../_specs/openapi.yml")

	t.Run("valid authenticated request returns profile with following=false", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			targetProfile := anyProfile(t)
			targetProfile.Following = false // ensure following is false
			followerID := nonEmptyString().Draw(t, "followerID")
			mux := http.NewServeMux()

			attachUnfollowUserFnToHttpHandler(mux, func(ctx context.Context, followerIDArg, followeeUsername string) (*domain.Profile, error) {
				assert.Equal(t, followerID, followerIDArg)
				assert.Equal(t, targetProfile.Username, followeeUsername)
				return &targetProfile, nil
			})

			requestParams := newUnfollowUserReqParams(targetProfile.Username)

			// when
			response := unfollowUser(mux, requestParams, &followerID)

			// then
			assert.Equal(t, http.StatusOK, response.StatusCode)

			respBody := parseResponseBody(tt, testutil.CloneResponse(tt, response))
			assert.Equal(t, targetProfile.Username, respBody.Profile.Username)
			assert.Equal(t, targetProfile.Bio, respBody.Profile.Bio)
			assert.Equal(t, targetProfile.Image, respBody.Profile.Image)
			assert.False(t, respBody.Profile.Following)

			testutil.EnsureRespMatchesSpecs(tt, tt.Context(), specsRouter, requestParams, response)
		})
	})

	t.Run("missing authentication returns 401", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			mux := http.NewServeMux()

			attachUnfollowUserFnToHttpHandler(mux, func(ctx context.Context, followerID, followeeUsername string) (*domain.Profile, error) {
				return nil, errors.New("should not be called")
			})

			requestParams := newUnfollowUserReqParams(nonEmptyString().Draw(t, "username"))

			// when
			response := unfollowUser(mux, requestParams, nil)

			// then
			assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		})
	})

	t.Run("user not found returns 404", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			mux := http.NewServeMux()

			attachUnfollowUserFnToHttpHandler(mux, func(ctx context.Context, followerID, followeeUsername string) (*domain.Profile, error) {
				return nil, usecase.ErrUserNotFound
			})

			requestParams := newUnfollowUserReqParams(nonEmptyString().Draw(t, "nonExistentUser"))

			// when
			followerID := nonEmptyString().Draw(t, "followerID")
			response := unfollowUser(mux, requestParams, &followerID)

			// then
			assert.Equal(t, http.StatusNotFound, response.StatusCode)

			testutil.EnsureRespMatchesSpecs(tt, tt.Context(), specsRouter, requestParams, response)
		})
	})

	t.Run("internal error returns 422", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			mux := http.NewServeMux()

			attachUnfollowUserFnToHttpHandler(mux, func(ctx context.Context, followerID, followeeUsername string) (*domain.Profile, error) {
				return nil, errors.New(nonEmptyString().Draw(t, "any unknown error"))
			})

			requestParams := newUnfollowUserReqParams(nonEmptyString().Draw(t, "username"))

			// when
			followerID := nonEmptyString().Draw(t, "followerID")
			response := unfollowUser(mux, requestParams, &followerID)

			// then
			assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)

			testutil.EnsureRespMatchesSpecs(tt, tt.Context(), specsRouter, requestParams, response)
		})
	})
}

// Helper to create an authenticated request for testing
func withAuthenticatedUser(req *http.Request, userInfo *auth.UserInfo) *http.Request {
	return auth.SetUserInContext(req, userInfo)
}

func anyProfile(t *rapid.T) domain.Profile {
	return domain.Profile{
		Username:  nonEmptyString().Draw(t, "username"),
		Bio:       nonEmptyString().Draw(t, "bio"),
		Image:     nonEmptyString().Draw(t, "image"),
		Following: rapid.Bool().Draw(t, "following"),
	}
}

func nonEmptyString() *rapid.Generator[string] {
	return rapid.StringOfN(
		rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")), // not too crazy string for now
		1, -1, -1,
	)
}

func getProfile(mux *http.ServeMux, apiReq testutil.APIRequest, as *string) *http.Response {
	req := httptest.NewRequest(apiReq.Method, apiReq.Path, nil)
	if as != nil {
		req = withAuthenticatedUser(req, &auth.UserInfo{ID: *as})
	}

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w.Result()
}

func followUser(mux *http.ServeMux, apiReq testutil.APIRequest, as *string) *http.Response {
	req := httptest.NewRequest(apiReq.Method, apiReq.Path, nil)
	if as != nil {
		req = withAuthenticatedUser(req, &auth.UserInfo{ID: *as})
	}

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w.Result()
}

func unfollowUser(mux *http.ServeMux, apiReq testutil.APIRequest, as *string) *http.Response {
	req := httptest.NewRequest(apiReq.Method, apiReq.Path, nil)
	if as != nil {
		req = withAuthenticatedUser(req, &auth.UserInfo{ID: *as})
	}

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w.Result()
}

func parseResponseBody(t *testing.T, resp *http.Response) profile.ProfileResponse {
	var profile profile.ProfileResponse
	err := json.NewDecoder(resp.Body).Decode(&profile)
	require.NoError(t, err)
	return profile
}

func newGetProfileReqParams(username string) testutil.APIRequest {
	return testutil.APIRequest{
		Method: "GET",
		Path:   "/api/profiles/" + username,
	}
}

func newFollowUserReqParams(username string) testutil.APIRequest {
	return testutil.APIRequest{
		Method: "POST",
		Path:   "/api/profiles/" + username + "/follow",
		Body:   "",
	}
}

func newUnfollowUserReqParams(username string) testutil.APIRequest {
	return testutil.APIRequest{
		Method: "DELETE",
		Path:   "/api/profiles/" + username + "/follow",
		Body:   "",
	}
}

func attachGetProfileFnToHttpHandler(mux *http.ServeMux, getProfileFn usecase.GetProfileFn) {
	profile.ProfileHandler{GetProfileFn: getProfileFn}.GetProfile(mux)
}

func attachFollowUserFnToHttpHandler(mux *http.ServeMux, followUserFn usecase.FollowUserFn) {
	profile.ProfileHandler{FollowUserFn: followUserFn}.FollowUser(mux)
}

func attachUnfollowUserFnToHttpHandler(mux *http.ServeMux, unfollowUserFn usecase.UnfollowUserFn) {
	profile.ProfileHandler{UnfollowUserFn: unfollowUserFn}.UnfollowUser(mux)
}
