package uc

import "github.com/err0r500/go-realworld-clean/domain"

func (i interactor) ProfileUpdateFollow(userName, followeeName string, follow bool) (*domain.User, error) {
	user, err := i.userRW.GetByName(userName)
	if err != nil {
		return nil, err
	}
	if user.Name != userName {
		return nil, errWrongUser
	}
	if user == nil {
		return nil, ErrNotFound
	}

	user.UpdateFollowees(followeeName, follow)

	if err := i.userRW.Save(*user); err != nil {
		return nil, err
	}

	return user, nil
}
