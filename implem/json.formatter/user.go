package formatter

import "github.com/err0r500/go-realworld-clean/domain"

type UserResp struct {
	Email    string  `json:"email"`
	Token    string  `json:"token"`
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

func NewUserResp(user domain.User, token string) UserResp {
	return UserResp{
		Email:    user.Email,
		Token:    token,
		Username: user.Name,
		Bio:      user.Bio,
		Image:    user.ImageLink,
	}
}
