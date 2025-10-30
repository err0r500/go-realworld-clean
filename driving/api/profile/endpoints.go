package profile

import (
	"errors"
	"log"
	"net/http"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/err0r500/go-realworld-clean/driving/api"
	"github.com/err0r500/go-realworld-clean/usecase"
)

// ProfileHandler handles profile-related HTTP endpoints
type ProfileHandler struct {
	usecase.GetProfileFn
	usecase.FollowUserFn
	usecase.UnfollowUserFn
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(profileUseCase *usecase.ProfileUseCases) *ProfileHandler {
	return &ProfileHandler{
		profileUseCase.GetProfile,
		profileUseCase.FollowUser,
		profileUseCase.UnfollowUser,
	}
}

// ProfileResponse wraps a profile for API responses
type ProfileResponse struct {
	Profile Profile `json:"profile"`
}

// Profile represents a user profile in API responses
type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

// GetProfile handles GET /api/profiles/{username}
// Authentication is optional - if authenticated, will show following status
func (h ProfileHandler) GetProfile(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/profiles/{username}", func(w http.ResponseWriter, r *http.Request) {
		username := r.PathValue("username")
		if username == "" {
			api.WriteError(w, http.StatusBadRequest, "Username is required")
			return
		}

		ctx := r.Context()

		// Try to get authenticated user (optional)
		var requestingUserID *string
		if userInfo, ok := auth.GetUserFromContext(r); ok {
			requestingUserID = &userInfo.ID
		}

		// Get profile from use case
		profile, err := h.GetProfileFn(ctx, requestingUserID, username)
		if err != nil {
			if errors.Is(err, usecase.ErrUserNotFound) {
				api.WriteError(w, http.StatusNotFound, "Profile not found")
			} else {
				log.Printf("Failed to get profile: %v", err)
				api.WriteError(w, http.StatusUnprocessableEntity, "Failed to get profile")
			}
			return
		}

		api.WriteJSON(w, http.StatusOK,
			ProfileResponse{
				Profile: Profile{
					Username:  profile.Username,
					Bio:       profile.Bio,
					Image:     profile.Image,
					Following: profile.Following,
				},
			})
	})
}

// FollowUser handles POST /api/profiles/{username}/follow
// Requires authentication
func (h ProfileHandler) FollowUser(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/profiles/{username}/follow",
		auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			username := r.PathValue("username")
			if username == "" {
				api.WriteError(w, http.StatusBadRequest, "Username is required")
				return
			}

			ctx := r.Context()

			// Get authenticated user
			userInfo, ok := auth.GetUserFromContext(r)
			if !ok {
				api.WriteError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// Follow user via use case
			profile, err := h.FollowUserFn(ctx, userInfo.ID, username)
			if err != nil {
				if errors.Is(err, usecase.ErrUserNotFound) {
					api.WriteError(w, http.StatusNotFound, "Profile not found")
				} else if errors.Is(err, usecase.ErrCannotFollowYourself) {
					api.WriteError(w, http.StatusUnprocessableEntity, "Cannot follow yourself")
				} else {
					log.Printf("Failed to follow user: %v", err)
					api.WriteError(w, http.StatusUnprocessableEntity, "Failed to follow user")
				}
				return
			}

			// Convert to response format
			response := ProfileResponse{
				Profile: Profile{
					Username:  profile.Username,
					Bio:       profile.Bio,
					Image:     profile.Image,
					Following: profile.Following,
				},
			}

			api.WriteJSON(w, http.StatusOK, response)
		}))
}

// UnfollowUser handles DELETE /api/profiles/{username}/follow
// Requires authentication
func (h ProfileHandler) UnfollowUser(mux *http.ServeMux) {
	mux.HandleFunc("DELETE /api/profiles/{username}/follow",
		auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			username := r.PathValue("username")
			if username == "" {
				api.WriteError(w, http.StatusBadRequest, "Username is required")
				return
			}

			ctx := r.Context()

			// Get authenticated user
			userInfo, ok := auth.GetUserFromContext(r)
			if !ok {
				api.WriteError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// Unfollow user via use case
			profile, err := h.UnfollowUserFn(ctx, userInfo.ID, username)
			if err != nil {
				if errors.Is(err, usecase.ErrUserNotFound) {
					api.WriteError(w, http.StatusNotFound, "Profile not found")
				} else {
					log.Printf("Failed to unfollow user: %v", err)
					api.WriteError(w, http.StatusUnprocessableEntity, "Failed to unfollow user")
				}
				return
			}

			// Convert to response format
			response := ProfileResponse{
				Profile: Profile{
					Username:  profile.Username,
					Bio:       profile.Bio,
					Image:     profile.Image,
					Following: profile.Following,
				},
			}

			api.WriteJSON(w, http.StatusOK, response)
		}))
}
