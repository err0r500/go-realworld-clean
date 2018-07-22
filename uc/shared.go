package uc

import "errors"

var (
	ErrUserNameAlreadyInUse = errors.New("this username is already in use")
	//ErrUserEmailAlreadyInUsed = errors.New("this email address is already in use")
	errWrongUser       = errors.New("woops, wrong user")
	errProfileNotFound = errors.New("profile not found")
	ErrUserNotFound    = errors.New("user not found")
	//errNilPointerWithoutError = errors.New("interface returned a nil pointer without error") // TODO: add interface name in error handling
)
