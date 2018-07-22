package formatter

import (
	"github.com/err0r500/go-realworld-clean/domain"
)

type Profile struct {
	Username  string  `json:"username"`
	Bio       *string `json:"bio,omitempty"`
	Picture   *string `json:"picture,omitempty"`
	Following bool    `json:"following"`
}

func NewProfileFromDomain(user domain.User, following bool) Profile {
	return Profile{
		Username:  user.Name,
		Bio:       user.Bio,
		Picture:   user.ImageLink,
		Following: following,
	}
}
