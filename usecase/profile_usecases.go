package usecase

import (
	"context"
	"fmt"

	"github.com/err0r500/go-realworld-clean/domain"
)

// ProfileUseCases handles profile-related business logic
type ProfileUseCases struct {
	profileRepo      ProfileRepository
	identityProvider IdentityProvider
}

// NewProfileUseCases creates a new profile use case
func NewProfileUseCases(profileRepo ProfileRepository, identityProvider IdentityProvider) *ProfileUseCases {
	return &ProfileUseCases{
		profileRepo:      profileRepo,
		identityProvider: identityProvider,
	}
}

type GetProfileFn = func(ctx context.Context, requestingUserID *string, targetUsername string) (*domain.Profile, error)
type UnfollowUserFn = func(ctx context.Context, followerID, followeeUsername string) (*domain.Profile, error)
type FollowUserFn = func(ctx context.Context, followerID, followeeUsername string) (*domain.Profile, error)

// GetProfile retrieves a user profile by username with follow status
// requestingUserID can be nil for unauthenticated requests
func (uc *ProfileUseCases) GetProfile(ctx context.Context, requestingUserID *string, targetUsername string) (*domain.Profile, error) {
	// Look up target user by username
	targetIdentity, err := uc.identityProvider.GetIdentityByUsername(ctx, targetUsername)
	if err != nil {
		return nil, err
	}

	// Create profile from identity
	profile := &domain.Profile{
		Username:  targetIdentity.Username,
		Bio:       targetIdentity.Bio,
		Image:     targetIdentity.Image,
		Following: false,
	}

	// Check follow status if user is authenticated
	if requestingUserID != nil && *requestingUserID != "" {
		isFollowing, err := uc.profileRepo.IsFollowing(*requestingUserID, targetIdentity.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check follow status: %w", err)
		}
		profile.Following = isFollowing
	}

	return profile, nil
}

// FollowUser creates a follow relationship
func (uc *ProfileUseCases) FollowUser(ctx context.Context, followerID, followeeUsername string) (*domain.Profile, error) {
	// Look up followee by username
	followeeIdentity, err := uc.identityProvider.GetIdentityByUsername(ctx, followeeUsername)
	if err != nil {
		return nil, err
	}

	// Check if trying to follow yourself
	if followerID == followeeIdentity.ID {
		return nil, ErrCannotFollowYourself
	}

	// Create follow relationship
	if err := uc.profileRepo.Follow(followerID, followeeIdentity.ID); err != nil {
		return nil, err
	}

	// Return updated profile
	return uc.GetProfile(ctx, &followerID, followeeUsername)
}

// UnfollowUser removes a follow relationship
func (uc *ProfileUseCases) UnfollowUser(ctx context.Context, followerID, followeeUsername string) (*domain.Profile, error) {
	// Look up followee by username
	followeeIdentity, err := uc.identityProvider.GetIdentityByUsername(ctx, followeeUsername)
	if err != nil {
		return nil, err
	}

	// Remove follow relationship
	if err := uc.profileRepo.Unfollow(followerID, followeeIdentity.ID); err != nil {
		return nil, err
	}

	// Return updated profile
	return uc.GetProfile(ctx, &followerID, followeeUsername)
}
