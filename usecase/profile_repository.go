package usecase

// ProfileRepository defines the interface for profile-related data operations
type ProfileRepository interface {
	// Follow creates a follow relationship between follower and followee
	Follow(followerID, followeeID string) error

	// Unfollow removes a follow relationship between follower and followee
	Unfollow(followerID, followeeID string) error

	// IsFollowing checks if follower is following followee
	IsFollowing(followerID, followeeID string) (bool, error)
}
