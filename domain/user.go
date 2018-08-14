package domain

import (
	"sort"
	"time"
)

// User represents a user account in the system
type User struct {
	Name      string
	Email     string
	Password  string
	Bio       *string
	ImageLink *string
	FollowIDs []string
	Favorites []Article
	CreatedAt time.Time
	UpdatedAt time.Time
}
type UserUpdatableProperty int

const (
	UserEmail UserUpdatableProperty = iota
	UserName
	UserBio
	UserImageLink
	UserPassword
)

func UpdateUser(initial *User, opts ...func(fields *User)) {
	for _, v := range opts {
		v(initial)
	}
}

func SetUserName(input *string) func(fields *User) {
	return func(initial *User) {
		if input != nil {
			initial.Name = *input
		}
	}
}

func SetUserEmail(input *string) func(fields *User) {
	return func(initial *User) {
		if input != nil {
			initial.Email = *input
		}
	}
}

// give empty string to delete it
func SetUserBio(input *string) func(fields *User) {
	return func(initial *User) {
		if input != nil {
			if *input == "" {
				initial.Bio = nil
				return
			}
			initial.Bio = input
		}
	}
}

// give empty string to delete it
func SetUserImageLink(input *string) func(fields *User) {
	return func(initial *User) {
		if input != nil {
			if *input == "" {
				initial.ImageLink = nil
				return
			}
			initial.ImageLink = input
		}
	}
}

func SetUserPassword(input *string) func(fields *User) {
	return func(initial *User) {
		if input != nil {
			initial.Password = *input
		}
	}
}

func (user User) Follows(userName string) bool {
	if user.FollowIDs == nil {
		return false
	}

	sort.Strings(user.FollowIDs)
	i := sort.SearchStrings(user.FollowIDs, userName)
	return i < len(user.FollowIDs) && user.FollowIDs[i] == userName
}

// UpdateFollowees will append or remove followee to current user according to follow param
func (user *User) UpdateFollowees(followeeName string, follow bool) {
	if follow {
		user.FollowIDs = append(user.FollowIDs, followeeName)
		return
	}

	for i := 0; i < len(user.FollowIDs); i++ {
		if user.FollowIDs[i] == followeeName {
			user.FollowIDs = append(user.FollowIDs[:i], user.FollowIDs[i+1:]...) // memory leak ? https://github.com/golang/go/wiki/SliceTricks
		}
	}
	if len(user.FollowIDs) == 0 {
		user.FollowIDs = nil
	}
}
