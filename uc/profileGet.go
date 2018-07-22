package uc

import (
	"github.com/err0r500/go-realworld-clean/domain"
)

func (i interactor) ProfileGet(userName string) (*domain.Profile, error) {
	user, err := i.userRW.GetByName(userName)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errProfileNotFound
	}
	if user.Name != userName {
		return nil, errWrongUser
	}

	return &domain.Profile{User: *user}, nil
}
