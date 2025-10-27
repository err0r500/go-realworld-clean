package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	ory "github.com/ory/kratos-client-go"
)

// User context key for storing authenticated user info
type contextKey string

const userContextKey contextKey = "user"

// UserInfo holds authenticated user information
type UserInfo struct {
	ID           string
	Email        string
	Username     string
	Bio          string
	Image        string
	SessionToken string // Kratos session token for self-service flows
}

// RequireAuth is a middleware that requires authentication
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, authenticated := GetUserFromContext(r)
		if !authenticated {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"errors":{"body":["Unauthorized"]}}`))
			return
		}
		next(w, r)
	}
}

// AuthMiddleware validates authentication via JWT Bearer token or session cookie
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Token "); ok {
			if user := validateJWTToken(r.Context(), token); user != nil {
				ctx := context.WithValue(r.Context(), userContextKey, user)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
		}

		// No valid authentication - continue without user context
		next.ServeHTTP(w, r)
	})
}

// validateJWTToken validates a JWT token and returns user info
func validateJWTToken(ctx context.Context, tokenString string) *UserInfo {
	// Step 1: Validate JWT structure, signature, and expiry
	claims, err := validateToken(tokenString)
	if err != nil {
		log.Printf("JWT validation failed: %v", err)
		return nil
	}

	// Step 2: Validate underlying Kratos session (checks for revocation)
	// This ensures the session hasn't been revoked even if JWT is valid
	session, err := validateSession(ctx, claims.KratosSessionToken)
	if err != nil {
		log.Printf("Kratos session validation failed: %v", err)
		return nil
	}

	return sessionToUserInfo(session, claims.KratosSessionToken)
}

// sessionToUserInfo converts a Kratos session to UserInfo
func sessionToUserInfo(session *ory.Session, sessionToken string) *UserInfo {
	// Check if session is active
	if !session.GetActive() {
		log.Printf("Session is not active")
		return nil
	}

	identity := session.GetIdentity()
	user := UserInfo{
		ID:           identity.GetId(),
		SessionToken: sessionToken,
	}

	// Extract traits
	if identity.Traits != nil {
		traits := identity.Traits.(map[string]any)
		if email, ok := traits["email"].(string); ok {
			user.Email = email
		}
		if username, ok := traits["username"].(string); ok {
			user.Username = username
		}
		if bio, ok := traits["bio"].(string); ok {
			user.Bio = bio
		}
		if image, ok := traits["image"].(string); ok {
			user.Image = image
		}
	}

	return &user
}

// GetUserFromContext retrieves user info from request context
func GetUserFromContext(r *http.Request) (*UserInfo, bool) {
	user, ok := r.Context().Value(userContextKey).(*UserInfo)
	if !ok {
		return nil, false
	}
	return user, true
}
