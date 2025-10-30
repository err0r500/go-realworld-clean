package usecase

import "errors"

// Sentinel errors for common use case failures
var (
	// ErrUserNotFound indicates the requested user does not exist
	ErrUserNotFound = errors.New("user not found")

	// ErrCannotFollowYourself indicates an attempt to follow one's own profile
	ErrCannotFollowYourself = errors.New("cannot follow yourself")
)
