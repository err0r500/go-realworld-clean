package uc

import "errors"

var (
	ErrAlreadyInUse = errors.New("this username is already in use")
	//ErrUserEmailAlreadyInUsed = errors.New("this email address is already in use")
	errWrongUser       = errors.New("woops, wrong user")
	errProfileNotFound = errors.New("profile not found")
	ErrNotFound        = errors.New("user not found")
	errArticleNotFound = errors.New("article not found")
)
