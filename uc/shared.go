package uc

import "errors"

var (
	ErrConflict = errors.New("this username is already in use")
	//ErrUserEmailAlreadyInUsed = errors.New("this email address is already in use")
	ErrUnauthorized    = errors.New("woops, wrong user")
	ErrProfileNotFound = errors.New("profile not found")
	ErrNotFound        = errors.New("user not found")
	errArticleNotFound = errors.New("article not found")
	errTechnical       = errors.New("a technical error happened")
	ErrTechnical       = errors.New("a technical error happened")
)
