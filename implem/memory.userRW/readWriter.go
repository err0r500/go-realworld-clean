package userRW

import (
	"context"
	"errors"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"time"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/err0r500/go-realworld-clean/uc"
	"github.com/spf13/viper"
)

type rw struct {
	store *sync.Map // map username:user
	//TODO : ADD a password hasher here
}

func New() uc.UserRW {
	rw := rw{
		store: &sync.Map{},
	}

	if viper.GetBool("populate") {
		rick := testData.User("rick")
		rw.Create(context.Background(), rick.Name, rick.Email, rick.Password)
	}

	return rw
}

func (rw rw) Create(ctx context.Context, username, email, password string) (*domain.User, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inmem_userrw:create")
	defer span.Finish()

	mayUser, ok := rw.GetByName(ctx, username)
	if !ok {
		return nil, false
	}
	if mayUser != nil {
		span.LogFields(log.Error(uc.ErrConflict))
	}

	rw.store.Store(username, domain.User{
		Name:      username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	u, ok := rw.GetByName(ctx, username)
	if !ok {
		return nil, false
	}
	return u, true
}

func (rw rw) GetByName(ctx context.Context, userName string) (*domain.User, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inmem_userrw:get_by_name")
	defer span.Finish()

	value, ok := rw.store.Load(userName)
	if !ok {
		return nil, true
	}

	user, ok := value.(domain.User)
	if !ok {
		span.LogFields(log.Error(uc.ErrTechnical))
		return nil, false
	}

	return &user, true
}

func (rw rw) GetByEmailAndPassword(ctx context.Context, email, password string) (*domain.User, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inmem_userrw:get_by_email_pass")
	defer span.Finish()

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

	if err != nil {
		return nil, false
	}

	return &foundUser, true
}

func (rw rw) Save(ctx context.Context, user domain.User) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inmem_userrw:save")
	defer span.Finish()

	mayUser, ok := rw.GetByName(ctx, user.Name)
	if !ok {
		span.LogFields(log.Error(uc.ErrTechnical))
		return false
	}
	if mayUser == nil {
		span.LogFields(log.Error(uc.ErrNotFound))
		return false
	}

	user.UpdatedAt = time.Now()
	rw.store.Store(user.Name, user)

	return true
}
