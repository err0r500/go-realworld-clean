package identityProvider

import (
	"context"
	"fmt"

	"github.com/err0r500/go-realworld-clean/usecase"
	ory "github.com/ory/kratos-client-go"
)

// KratosIdentityProvider implements IdentityProvider using Ory Kratos
type KratosIdentityProvider struct {
	adminClient *ory.APIClient
}

// NewKratosIdentityProvider creates a new Kratos identity provider
func NewKratosIdentityProvider(adminClient *ory.APIClient) *KratosIdentityProvider {
	return &KratosIdentityProvider{
		adminClient: adminClient,
	}
}

// GetUserByUsername retrieves a user identity by username
// Note: Kratos doesn't have a direct way to search by username trait, so we list all
// identities and filter by username. NOT PRODUCTION READY
func (k *KratosIdentityProvider) GetIdentityByUsername(ctx context.Context, username string) (*usecase.UserIdentity, error) {
	identities, _, err := k.adminClient.IdentityAPI.
		ListIdentities(ctx).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to list identities: %w", err)
	}

	// Filter by username trait
	for _, identity := range identities {
		if identity.Traits != nil {
			traits := identity.Traits.(map[string]any)
			if usernameVal, ok := traits["username"].(string); ok && usernameVal == username {
				return k.kratosIdentityToUserIdentity(&identity), nil
			}
		}
	}

	return nil, usecase.ErrUserNotFound
}

// GetUserByID retrieves a user identity by ID
func (k *KratosIdentityProvider) GetIdentityByID(ctx context.Context, id string) (*usecase.UserIdentity, error) {
	identity, _, err := k.adminClient.IdentityAPI.
		GetIdentity(ctx, id).
		Execute()

	if err != nil {
		return nil, usecase.ErrUserNotFound
	}

	return k.kratosIdentityToUserIdentity(identity), nil
}

// kratosIdentityToUserIdentity converts a Kratos Identity to UserIdentity
func (k *KratosIdentityProvider) kratosIdentityToUserIdentity(identity *ory.Identity) *usecase.UserIdentity {
	userIdentity := &usecase.UserIdentity{
		ID: identity.GetId(),
	}

	// Extract traits
	if identity.Traits != nil {
		traits := identity.Traits.(map[string]any)
		if username, ok := traits["username"].(string); ok {
			userIdentity.Username = username
		}
		if email, ok := traits["email"].(string); ok {
			userIdentity.Email = email
		}
		if bio, ok := traits["bio"].(string); ok {
			userIdentity.Bio = bio
		}
		if image, ok := traits["image"].(string); ok {
			userIdentity.Image = image
		}
	}

	return userIdentity
}
