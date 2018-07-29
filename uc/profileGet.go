package uc

import (
	"github.com/err0r500/go-realworld-clean/domain"
)

func (i interactor) ProfileGet(requestingUserName, userName string) (*domain.User, bool, error) {
	user, err := i.userRW.GetByName(userName)
	if err != nil {
		return nil, false, err
	}
	if user == nil {
		return nil, false, errProfileNotFound
	}

	if requestingUserName == "" {
		return user, false, nil
	}

	reqUser, err := i.userRW.GetByName(requestingUserName)
	if err != nil {
		return nil, false, err
	}
	if user == nil {
		return nil, false, errProfileNotFound
	}

	return user, reqUser.Follows(userName), nil
}
