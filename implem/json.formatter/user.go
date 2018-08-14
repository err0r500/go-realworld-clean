package formatter

import "github.com/err0r500/go-realworld-clean/domain"

type UserResp struct {
	Email     string `json:"email"`
	Token     string `json:"token"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func NewUserResp(user domain.User, token string) UserResp {
	var bio, image string
	if user.Bio != nil {
		bio = *user.Bio
	}
	if user.ImageLink != nil {
		image = *user.ImageLink
	}
	return UserResp{
		Email:     user.Email,
		Token:     token,
		Username:  user.Name,
		Bio:       bio,
		Image:     image,
		CreatedAt: user.CreatedAt.UTC().Format(dateLayout),
		UpdatedAt: user.UpdatedAt.UTC().Format(dateLayout),
	}
}
