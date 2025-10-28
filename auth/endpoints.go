package auth

import (
	"context"
	"encoding/json"
	"log"
	"maps"
	"net/http"
	"strings"

	"github.com/err0r500/go-realworld-clean/driving/api"
	ory "github.com/ory/kratos-client-go"
)

type AuthHandler struct {
	kratosClient *ory.APIClient
}

func NewAuthHandler(kratosURL string) *AuthHandler {
	config := ory.NewConfiguration()
	config.Servers = ory.ServerConfigurations{
		{URL: kratosURL},
	}

	return &AuthHandler{
		kratosClient: ory.NewAPIClient(config),
	}
}

type RegisterUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterReqBody struct {
	User RegisterUser `json:"user"`
}

func (auth AuthHandler) RegisterUser(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/users",
		func(w http.ResponseWriter, r *http.Request) {
			var req RegisterReqBody
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				api.WriteError(w, http.StatusBadRequest, "Invalid request body")
				return
			}

			// Validate required fields
			if req.User.Email == "" || req.User.Password == "" || req.User.Username == "" {
				api.WriteError(w, http.StatusUnprocessableEntity, "Email, password, and username are required")
				return
			}

			ctx := r.Context()

			// Step 1: Initialize registration flow
			flow, _, err := auth.kratosClient.FrontendAPI.CreateNativeRegistrationFlow(ctx).Execute()
			if err != nil {
				log.Printf("Failed to create registration flow: %v", err)
				api.WriteError(w, http.StatusInternalServerError, "Failed to initialize registration")
				return
			}

			// Step 2: Submit registration with user data
			updateBody := ory.UpdateRegistrationFlowWithPasswordMethodAsUpdateRegistrationFlowBody(
				&ory.UpdateRegistrationFlowWithPasswordMethod{
					Method:   "password",
					Password: req.User.Password,
					Traits: map[string]any{
						"email":    req.User.Email,
						"username": req.User.Username,
						"bio":      "",
						"image":    "",
					},
				},
			)

			result, resp, err := auth.kratosClient.FrontendAPI.UpdateRegistrationFlow(ctx).
				Flow(flow.GetId()).
				UpdateRegistrationFlowBody(updateBody).
				Execute()

			if err != nil {
				log.Printf("Registration failed: %v", err)
				if resp != nil && resp.StatusCode == 400 {
					api.WriteError(w, http.StatusUnprocessableEntity, "Registration failed - email may already be in use")
				} else {
					api.WriteError(w, http.StatusInternalServerError, "Registration failed")
				}
				return
			}

			session := result.GetSession()
			identity := session.GetIdentity()

			// Step 5: Generate JWT token with session info
			token, err := GenerateToken(
				result.GetSessionToken(),
				identity.GetId(),
				session.GetExpiresAt(),
			)
			if err != nil {
				log.Printf("Failed to generate token: %v", err)
				api.WriteError(w, http.StatusInternalServerError, "Failed to generate ntication token")
				return
			}

			api.WriteJSON(w, http.StatusCreated, UserResponse{User: IdentityToUser(&identity, token)})
		})
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginReqBody struct {
	User LoginUser `json:"user"`
}

func (auth AuthHandler) Login(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/users/login",
		func(w http.ResponseWriter, r *http.Request) {
			var req LoginReqBody
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				api.WriteError(w, http.StatusBadRequest, "Invalid request body")
				return
			}

			// Validate required fields
			if req.User.Email == "" || req.User.Password == "" {
				api.WriteError(w, http.StatusUnprocessableEntity, "Email and password are required")
				return
			}

			ctx := r.Context()

			// Step 1: Initialize login flow
			flow, _, err := auth.kratosClient.FrontendAPI.CreateNativeLoginFlow(ctx).Execute()
			if err != nil {
				log.Printf("Failed to create login flow: %v", err)
				api.WriteError(w, http.StatusInternalServerError, "Failed to initialize login")
				return
			}

			// Step 2: Submit login credentials
			result, resp, err := auth.kratosClient.FrontendAPI.UpdateLoginFlow(ctx).
				Flow(flow.GetId()).
				UpdateLoginFlowBody(
					ory.UpdateLoginFlowWithPasswordMethodAsUpdateLoginFlowBody(
						&ory.UpdateLoginFlowWithPasswordMethod{
							Method:     "password",
							Identifier: req.User.Email,
							Password:   req.User.Password,
						},
					)).
				Execute()

			if err != nil {
				if resp == nil {
					api.WriteError(w, http.StatusInternalServerError, "Login failed")
					return
				}

				switch resp.StatusCode {
				case 400:
					api.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
				case 410:
					api.WriteError(w, http.StatusUnprocessableEntity, "Login flow expired - please try again")

				default:
					api.WriteError(w, http.StatusInternalServerError, "Login failed")
				}
				return
			}

			session := result.GetSession()
			identity := session.GetIdentity()

			// Step 4: Generate JWT token with session info
			token, err := GenerateToken(
				result.GetSessionToken(),
				identity.GetId(),
				session.GetExpiresAt(),
			)
			if err != nil {
				api.WriteError(w, http.StatusInternalServerError, "Failed to generate ntication token")
				return
			}

			api.WriteJSON(w, http.StatusOK, UserResponse{User: IdentityToUser(&identity, token)})
		})
}

func (auth AuthHandler) GetCurrentUser(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/user",
		RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			userInfo, _ := GetUserFromContext(r)
			token, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Token ")
			api.WriteJSON(w, http.StatusOK, UserResponse{User: User{
				Email:    userInfo.Email,
				Username: userInfo.Username,
				Bio:      userInfo.Bio,
				Image:    userInfo.Image,
				Token:    token,
			}})
		}))
}

type UpdateUserFields struct {
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	Image    *string `json:"image,omitempty"`
}

type UpdateUserReqBody struct {
	User UpdateUserFields `json:"user"`
}

func (auth AuthHandler) UpdateUser(mux *http.ServeMux) {
	mux.HandleFunc("PUT /api/user",
		RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			var req UpdateUserReqBody
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				api.WriteError(w, http.StatusBadRequest, "Invalid request body")
				return
			}

			userInfo, _ := GetUserFromContext(r)

			// Fetch current identity
			identity, _, err := auth.kratosClient.IdentityAPI.GetIdentity(r.Context(), userInfo.ID).Execute()
			if err != nil {
				log.Printf("Failed to fetch identity: %v", err)
				api.WriteError(w, http.StatusInternalServerError, "Failed to fetch user data")
				return
			}

			// Build updated traits - merge with existing
			traits := make(map[string]any)
			if identity.Traits != nil {
				existingTraits := identity.Traits.(map[string]any)
				maps.Copy(traits, existingTraits)
			}

			// Update only provided fields
			if req.User.Email != nil {
				traits["email"] = *req.User.Email
			}
			if req.User.Username != nil {
				traits["username"] = *req.User.Username
			}
			if req.User.Bio != nil {
				traits["bio"] = *req.User.Bio
			}
			if req.User.Image != nil {
				traits["image"] = *req.User.Image
			}

			// Handle password update separately if provided
			if req.User.Password != nil && *req.User.Password != "" {
				if err := auth.updatePassword(r.Context(), userInfo.SessionToken, *req.User.Password); err != nil {
					log.Printf("Failed to update password: %v", err)
					api.WriteError(w, http.StatusInternalServerError, "Failed to update password")
					return
				}
			}

			updatedIdentity, _, err := auth.kratosClient.IdentityAPI.UpdateIdentity(r.Context(), userInfo.ID).
				UpdateIdentityBody(
					ory.UpdateIdentityBody{
						SchemaId: identity.SchemaId,
						State:    *identity.State,
						Traits:   traits,
					}).
				Execute()

			if err != nil {
				log.Printf("Failed to update identity: %v", err)
				api.WriteError(w, http.StatusInternalServerError, "Failed to update user")
				return
			}

			// Return updated user
			api.WriteJSON(w, http.StatusOK, UserResponse{User: IdentityToUser(updatedIdentity, "")})
		}))
}

// updatePassword updates the user's password using Kratos settings flow
func (auth AuthHandler) updatePassword(ctx context.Context, sessionToken, newPassword string) error {
	// Step 1: Initialize settings flow
	flow, _, err := auth.kratosClient.FrontendAPI.CreateNativeSettingsFlow(ctx).
		XSessionToken(sessionToken).
		Execute()
	if err != nil {
		return err
	}

	// Step 2: Submit password update
	updateBody := ory.UpdateSettingsFlowWithPasswordMethodAsUpdateSettingsFlowBody(
		&ory.UpdateSettingsFlowWithPasswordMethod{
			Method:   "password",
			Password: newPassword,
		},
	)

	_, _, err = auth.kratosClient.FrontendAPI.UpdateSettingsFlow(ctx).
		Flow(flow.GetId()).
		XSessionToken(sessionToken).
		UpdateSettingsFlowBody(updateBody).
		Execute()

	return err
}
