package usecase_test

import (
	"context"
	"log"
	"testing"

	th "github.com/err0r500/go-realworld-clean/test_helpers"
	"github.com/err0r500/go-realworld-clean/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestGetProfile(t *testing.T) {
	t.Run("get profile without authentication returns profile, following always false", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			identity := anyIdentity(t)
			uc := usecase.NewProfileUseCases(mockProfileRepo{},
				mockIdentityProvider{
					getByUsernameFn: getByUsername{
						expectedUsername: identity.Username,
						identity:         &identity,
					},
				})

			// when
			profile, err := uc.GetProfile(t.Context(), nil, identity.Username)

			// then
			assert.NoError(t, err)
			assert.Equal(t, identity.Username, profile.Username)
			assert.Equal(t, identity.Bio, profile.Bio)
			assert.Equal(t, identity.Image, profile.Image)
			assert.False(t, profile.Following)
		})
	})

	t.Run("get profile with authentication - returns following status", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			requestingUserID := th.NonEmptyString().Draw(t, "requestingUserID")
			requestingUserFollows := rapid.Bool().Draw(t, "follows")
			followeeIdentity := anyIdentity(t)

			uc := usecase.NewProfileUseCases(
				mockProfileRepo{
					isFollowingFn: isfollowing{
						followerID: requestingUserID,
						followeeID: followeeIdentity.ID,
						resp:       requestingUserFollows,
					},
				},
				mockIdentityProvider{
					getByUsernameFn: getByUsername{
						expectedUsername: followeeIdentity.Username,
						identity:         &followeeIdentity,
					},
				},
			)

			// when
			profile, err := uc.GetProfile(t.Context(), &requestingUserID, followeeIdentity.Username)

			// then
			require.NoError(t, err)
			assert.Equal(t, requestingUserFollows, profile.Following)
		})
	})

	t.Run("user not found", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			username := th.NonEmptyString().Draw(t, "username")
			uc := usecase.NewProfileUseCases(
				mockProfileRepo{},
				mockIdentityProvider{
					getByUsernameFn: getByUsername{
						expectedUsername: username,
						err:              usecase.ErrUserNotFound,
					},
				},
			)

			// when
			profile, err := uc.GetProfile(t.Context(), nil, username)

			// then
			assert.Nil(t, profile)
			assert.ErrorIs(t, err, usecase.ErrUserNotFound)
		})
	})

	t.Run("repository error when checking follow status", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			requestingUserID := th.NonEmptyString().Draw(t, "requestingUserID")
			followeeIdentity := anyIdentity(t)

			uc := usecase.NewProfileUseCases(
				mockProfileRepo{
					isFollowingFn: isfollowing{
						followerID: requestingUserID,
						followeeID: followeeIdentity.ID,
						err:        assert.AnError,
					},
				},
				mockIdentityProvider{
					getByUsernameFn: getByUsername{
						expectedUsername: followeeIdentity.Username,
						identity:         &followeeIdentity,
					},
				},
			)

			// when
			profile, err := uc.GetProfile(t.Context(), &requestingUserID, followeeIdentity.Username)

			// then
			assert.Nil(t, profile)
			assert.Error(t, err)
		})
	})
}

func TestFollowUser(t *testing.T) {
	t.Run("successfully follow user", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			followerID := th.NonEmptyString().Draw(t, "followerID") + "-follower"
			followeeIdentity := anyIdentity(t)

			uc := usecase.NewProfileUseCases(
				mockProfileRepo{
					followFunc: follow{
						followerID: followerID,
						followeeID: followeeIdentity.ID,
					},
					isFollowingFn: isfollowing{
						followerID: followerID,
						followeeID: followeeIdentity.ID,
						resp:       true,
					},
				},
				mockIdentityProvider{
					getByUsernameFn: getByUsername{
						expectedUsername: followeeIdentity.Username,
						identity:         &followeeIdentity,
					},
				},
			)

			// when
			profile, err := uc.FollowUser(t.Context(), followerID, followeeIdentity.Username)

			// then
			require.NoError(t, err)
			assert.Equal(t, followeeIdentity.Username, profile.Username)
			assert.Equal(t, followeeIdentity.Bio, profile.Bio)
			assert.Equal(t, followeeIdentity.Image, profile.Image)
			assert.True(t, profile.Following)
		})
	})

	t.Run("cannot follow yourself", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			identity := anyIdentity(t)

			uc := usecase.NewProfileUseCases(
				mockProfileRepo{},
				mockIdentityProvider{
					getByUsernameFn: getByUsername{
						expectedUsername: identity.Username,
						identity:         &identity,
					},
				},
			)

			// when
			profile, err := uc.FollowUser(t.Context(), identity.ID, identity.Username)

			// then
			assert.Nil(t, profile)
			assert.ErrorIs(t, err, usecase.ErrCannotFollowYourself)
		})
	})

	t.Run("followee user not found", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			followerID := th.NonEmptyString().Draw(t, "followerID")
			username := th.NonEmptyString().Draw(t, "username")

			uc := usecase.NewProfileUseCases(
				mockProfileRepo{},
				mockIdentityProvider{
					getByUsernameFn: getByUsername{
						expectedUsername: username,
						err:              usecase.ErrUserNotFound,
					},
				},
			)

			// when
			profile, err := uc.FollowUser(t.Context(), followerID, username)

			// then
			assert.Nil(t, profile)
			assert.ErrorIs(t, err, usecase.ErrUserNotFound)
		})
	})

	t.Run("repository error when following", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			followerID := th.NonEmptyString().Draw(t, "followerID") + "-follower"
			followeeIdentity := anyIdentity(t)

			uc := usecase.NewProfileUseCases(
				mockProfileRepo{
					followFunc: follow{
						followerID: followerID,
						followeeID: followeeIdentity.ID,
						err:        assert.AnError,
					},
				},
				mockIdentityProvider{
					getByUsernameFn: getByUsername{
						expectedUsername: followeeIdentity.Username,
						identity:         &followeeIdentity,
					},
				},
			)

			// when
			profile, err := uc.FollowUser(t.Context(), followerID, followeeIdentity.Username)

			// then
			assert.Nil(t, profile)
			assert.Error(t, err)
		})
	})
}

func TestUnfollowUser(t *testing.T) {
	t.Run("successfully unfollow user", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			followerID := th.NonEmptyString().Draw(t, "followerID")
			followeeIdentity := anyIdentity(t)

			uc := usecase.NewProfileUseCases(
				mockProfileRepo{
					unfollowFunc: follow{
						followerID: followerID,
						followeeID: followeeIdentity.ID,
					},
					isFollowingFn: isfollowing{
						followerID: followerID,
						followeeID: followeeIdentity.ID,
						resp:       false,
					},
				},
				mockIdentityProvider{
					getByUsernameFn: getByUsername{
						expectedUsername: followeeIdentity.Username,
						identity:         &followeeIdentity,
					},
				},
			)

			// when
			profile, err := uc.UnfollowUser(t.Context(), followerID, followeeIdentity.Username)

			// then
			require.NoError(t, err)
			assert.Equal(t, followeeIdentity.Username, profile.Username)
			assert.False(t, profile.Following)
		})
	})

	t.Run("unfollow when not following - no error", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			followerID := th.NonEmptyString().Draw(t, "followerID")
			followeeIdentity := anyIdentity(t)

			uc := usecase.NewProfileUseCases(
				mockProfileRepo{
					unfollowFunc: follow{
						followerID: followerID,
						followeeID: followeeIdentity.ID,
					},
					isFollowingFn: isfollowing{
						followerID: followerID,
						followeeID: followeeIdentity.ID,
						resp:       false,
					},
				},
				mockIdentityProvider{
					getByUsernameFn: getByUsername{
						expectedUsername: followeeIdentity.Username,
						identity:         &followeeIdentity,
					},
				},
			)

			// when
			profile, err := uc.UnfollowUser(t.Context(), followerID, followeeIdentity.Username)

			// then
			require.NoError(t, err)
			assert.Equal(t, followeeIdentity.Username, profile.Username)
			assert.False(t, profile.Following)
		})
	})

	t.Run("user not found", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			followerID := th.NonEmptyString().Draw(t, "followerID")
			username := th.NonEmptyString().Draw(t, "username")

			uc := usecase.NewProfileUseCases(
				mockProfileRepo{},
				mockIdentityProvider{
					getByUsernameFn: getByUsername{
						expectedUsername: username,
						err:              usecase.ErrUserNotFound,
					},
				},
			)

			// when
			profile, err := uc.UnfollowUser(t.Context(), followerID, username)

			// then
			assert.Nil(t, profile)
			assert.ErrorIs(t, err, usecase.ErrUserNotFound)
		})
	})

	t.Run("repository error when unfollowing", func(tt *testing.T) {
		rapid.Check(tt, func(t *rapid.T) {
			// given
			followerID := th.NonEmptyString().Draw(t, "followerID")
			followeeIdentity := anyIdentity(t)

			uc := usecase.NewProfileUseCases(
				mockProfileRepo{
					unfollowFunc: follow{
						followerID: followerID,
						followeeID: followeeIdentity.ID,
						err:        assert.AnError,
					},
				},
				mockIdentityProvider{
					getByUsernameFn: getByUsername{
						expectedUsername: followeeIdentity.Username,
						identity:         &followeeIdentity,
					},
				},
			)

			// when
			profile, err := uc.UnfollowUser(t.Context(), followerID, followeeIdentity.Username)

			// then
			assert.Nil(t, profile)
			assert.Error(t, err)
		})
	})
}

type follow struct {
	followerID string
	followeeID string
	err        error
}

type isfollowing struct {
	followerID string
	followeeID string
	resp       bool
	err        error
}

// In-memory mock for ProfileRepository
type mockProfileRepo struct {
	followFunc    follow
	unfollowFunc  follow
	isFollowingFn isfollowing
}

func (m mockProfileRepo) Follow(followerID, followeeID string) error {
	if followerID != m.followFunc.followerID {
		log.Fatalf("followerID, expected: %s, got %s", m.followFunc.followerID, followerID)
	}

	if followeeID != m.followFunc.followeeID {
		log.Fatalf("followeeID, expected: %s, got %s", m.followFunc.followeeID, followeeID)
	}
	return m.followFunc.err
}

func (m mockProfileRepo) Unfollow(followerID, followeeID string) error {
	if followerID != m.unfollowFunc.followerID {
		log.Fatalf("followerID, expected: %s, got %s", m.unfollowFunc.followerID, followerID)
	}

	if followeeID != m.unfollowFunc.followeeID {
		log.Fatalf("followeeID, expected: %s, got %s", m.unfollowFunc.followeeID, followeeID)
	}
	return m.unfollowFunc.err
}

func (m mockProfileRepo) IsFollowing(followerID, followeeID string) (bool, error) {
	if followerID != m.isFollowingFn.followerID {
		log.Fatalf("followerID, expected: %s, got %s", m.unfollowFunc.followerID, followerID)
	}

	if followeeID != m.isFollowingFn.followeeID {
		log.Fatalf("followeeID, expected: %s, got %s", m.unfollowFunc.followeeID, followeeID)
	}
	return m.isFollowingFn.resp, m.isFollowingFn.err
}

type getByUsername struct {
	expectedUsername string
	identity         *usecase.UserIdentity
	err              error
}

// In-memory mock for IdentityProvider
type mockIdentityProvider struct {
	getByUsername   func(username string) (*usecase.UserIdentity, error)
	getByID         func(id string) (*usecase.UserIdentity, error)
	getByUsernameFn getByUsername
}

func (m mockIdentityProvider) GetIdentityByUsername(ctx context.Context, username string) (*usecase.UserIdentity, error) {
	if m.getByUsernameFn.expectedUsername != username {
		log.Fatalf("GetIdentifyByUsername called with wrong username: expected %s, got %s", m.getByUsernameFn.expectedUsername, username)
	}
	return m.getByUsernameFn.identity, m.getByUsernameFn.err
}

func (m mockIdentityProvider) GetIdentityByID(ctx context.Context, id string) (*usecase.UserIdentity, error) {
	return m.getByID(id)
}

func anyIdentity(t *rapid.T) usecase.UserIdentity {
	return usecase.UserIdentity{
		ID:       th.NonEmptyString().Draw(t, "ID"),
		Username: th.NonEmptyString().Draw(t, "username"),
		Email:    th.NonEmptyString().Draw(t, "email"),
		Bio:      th.NonEmptyString().Draw(t, "bio"),
		Image:    th.NonEmptyString().Draw(t, "image"),
	}
}
