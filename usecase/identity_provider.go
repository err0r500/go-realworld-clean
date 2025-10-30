package usecase

import "context"

// IdentityProvider defines the interface for identity and user management operations
type IdentityProvider interface {
	// GetIdentityByUsername retrieves a user identity by username
	GetIdentityByUsername(ctx context.Context, username string) (*UserIdentity, error)

	// GetIdentityByID retrieves a user identity by ID
	GetIdentityByID(ctx context.Context, id string) (*UserIdentity, error)
}

// UserIdentity represents a user's identity information
type UserIdentity struct {
	ID       string
	Username string
	Email    string
	Bio      string
	Image    string
}
