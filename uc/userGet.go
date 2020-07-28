package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) UserGet(ctx context.Context, userName string) (*domain.User, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:user_get")
	defer span.Finish()

	user, ok := i.userRW.GetByName(ctx, userName)
	if !ok {
		return nil, "", ErrTechnical
	}
	if user == nil {
		return nil, "", ErrNotFound
	}
	if user.Name != userName {
		return nil, "", ErrUnauthorized
	}

	token, ok := i.authHandler.GenUserToken(ctx, userName)
	if !ok {
		return nil, "", ErrTechnical
	}

	return user, token, nil
}
