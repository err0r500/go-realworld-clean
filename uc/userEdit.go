package uc

import (
	"github.com/err0r500/go-realworld-clean/domain"
)

//UserEdit(userName string, newUser map[UserUpdatableProperty]*string) (user *domain.User, err error)
func (i interactor) UserEdit(userName string, fieldsToUpdate map[domain.UserUpdatableProperty]*string) (*domain.User, string, error) {
	user, err := i.userRW.GetByName(userName)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", ErrNotFound
	}
	if user.Name != userName {
		return nil, "", errWrongUser
	}

	domain.UpdateUser(user,
		domain.SetUserName(fieldsToUpdate[domain.UserName]),
		domain.SetUserEmail(fieldsToUpdate[domain.UserEmail]),
		domain.SetUserBio(fieldsToUpdate[domain.UserBio]),
		domain.SetUserImageLink(fieldsToUpdate[domain.UserImageLink]),
		domain.SetUserPassword(fieldsToUpdate[domain.UserPassword]),
	)

	if err := i.userValidator.CheckUser(*user); err != nil {
		return nil, "", err
	}

	if err := i.userRW.Save(*user); err != nil {
		return nil, "", err
	}

	token, err := i.authHandler.GenUserToken(user.Name)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
