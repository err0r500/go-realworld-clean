package userRW

import (
	"errors"
	"sync"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/uc"
)

type rw struct {
	store *sync.Map // map username:user
	//TODO : ADD a password hasher here
}

func New() uc.UserRW {
	return rw{
		store: &sync.Map{},
	}
}

func (rw rw) Create(username, email, password string) (*domain.User, error) {
	if _, err := rw.GetByName(username); err == nil {
		return nil, uc.ErrAlreadyInUse
	}

	rw.store.Store(username, domain.User{
		Name:     username,
		Email:    email,
		Password: password,
	})

	return rw.GetByName(username)
}

func (rw rw) GetByName(userName string) (*domain.User, error) {
	value, ok := rw.store.Load(userName)
	if !ok {
		return nil, uc.ErrNotFound
	}

	user, ok := value.(domain.User)
	if !ok {
		return nil, errors.New("not a user stored at key")
	}

	return &user, nil
}

func (rw rw) GetByEmailAndPassword(email, password string) (*domain.User, error) {
	var err error
	var foundUser domain.User

	rw.store.Range(func(key, value interface{}) bool {
		user, ok := value.(domain.User)
		if !ok {
			err = errors.New("failed to assert to domain.User")
			return false
		}

		if user.Email == email && user.Password == password {
			foundUser = user
			return false // stop range
		}

		return true // keep iterating
	})

	return &foundUser, err
}

func (rw rw) Save(user domain.User) error {
	if user, _ := rw.GetByName(user.Name); user == nil {
		return uc.ErrNotFound
	}

	rw.store.Store(user.Name, user)

	return nil
}
