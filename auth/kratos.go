package auth

import (
	"context"
	"log"
	"os"

	ory "github.com/ory/kratos-client-go"
)

var (
	kratosClient *ory.APIClient
	kratosURL    string
)

// InitKratos initializes the Kratos client
func InitKratos() {
	kratosURL = getEnvOrDefault("KRATOS_PUBLIC_URL", "http://localhost:4433")

	config := ory.NewConfiguration()
	config.Servers = ory.ServerConfigurations{
		{URL: kratosURL},
	}

	kratosClient = ory.NewAPIClient(config)
	log.Printf("Kratos client initialized with URL: %s", kratosURL)
}

// GetKratosClient returns the initialized Kratos client
func GetKratosClient() *ory.APIClient {
	if kratosClient == nil {
		InitKratos()
	}
	return kratosClient
}

// validateSession validates a Kratos session and returns the session data
func validateSession(ctx context.Context, sessionToken string) (*ory.Session, error) {
	client := GetKratosClient()

	session, _, err := client.FrontendAPI.ToSession(ctx).
		XSessionToken(sessionToken).
		Execute()

	return session, err
}

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

func getEnvOrDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}
