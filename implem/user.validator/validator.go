package validator

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/uc"
)

type userValidator struct{}

func New() uc.UserValidator {
	return userValidator{}
}

func (userValidator) CheckUser(user domain.User) error {
	if ok := govalidator.IsEmail(user.Email); !ok {
		return errors.New("invalid email")
	}

	return nil
}
