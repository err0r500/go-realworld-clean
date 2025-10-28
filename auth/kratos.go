package auth

import (
	"context"
	"os"

	ory "github.com/ory/kratos-client-go"
)

// IdentityToUser converts a Kratos Identity to a RealWorld User
func IdentityToUser(identity *ory.Identity, token string) User {
	user := User{
		Token: token,
	}

	// Extract traits from identity
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

	return user
}

// validateSession validates a Kratos session and returns the session data
func validateSession(ctx context.Context, sessionToken string, client *ory.APIClient) (*ory.Session, error) {
	session, _, err := client.FrontendAPI.ToSession(ctx).
		XSessionToken(sessionToken).
		Execute()

	return session, err
}

func getEnvOrDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}
