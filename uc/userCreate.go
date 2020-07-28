package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) UserCreate(ctx context.Context, username, email, password string) (*domain.User, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:user_create")
	defer span.Finish()

	user, ok := i.userRW.Create(ctx, username, email, password)
	if !ok {
		return nil, "", ErrTechnical
	}

	token, ok := i.authHandler.GenUserToken(ctx, username)
	if !ok {
		return nil, "", ErrTechnical
	}

	return user, token, nil
}
