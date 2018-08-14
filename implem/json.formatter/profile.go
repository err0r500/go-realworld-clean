package formatter

import (
	"github.com/err0r500/go-realworld-clean/domain"
)

type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Picture   string `json:"image"`
	Following bool   `json:"following"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func NewProfileFromDomain(user domain.User, following bool) Profile {
	var bio, image string
	if user.Bio != nil {
		bio = *user.Bio
	}
	if user.ImageLink != nil {
		image = *user.ImageLink
	}

	return Profile{
		Username:  user.Name,
		Bio:       bio,
		Picture:   image,
		Following: following,
		CreatedAt: user.CreatedAt.UTC().Format(dateLayout),
		UpdatedAt: user.UpdatedAt.UTC().Format(dateLayout),
	}
}
