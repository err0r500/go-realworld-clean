package uc

import (
	"github.com/err0r500/go-realworld-clean/domain"
)

func (i interactor) UserCreate(username, email, password string) (*domain.User, string, error) {
	user, err := i.userRW.Create(username, email, password)
	if err != nil {
		return nil, "", err
	}

	token, err := i.authHandler.GenUserToken(username)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
