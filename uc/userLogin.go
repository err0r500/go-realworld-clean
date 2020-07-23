package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) UserLogin(ctx context.Context, email, password string) (*domain.User, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:user_login")
	defer span.Finish()

	mayUser, ok := i.userRW.GetByEmailAndPassword(ctx, email, password)
	if !ok {
		return nil, "", ErrTechnical
	}
	if mayUser == nil {
		return nil, "", ErrNotFound
	}

	token, err := i.authHandler.GenUserToken(mayUser.Name)
	if err != nil {
		return nil, "", err
	}

	return mayUser, token, nil
}
