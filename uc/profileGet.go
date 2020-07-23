package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) ProfileGet(ctx context.Context, requestingUserName, userName string) (*domain.User, bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:profile_get")
	defer span.Finish()

	user, ok := i.userRW.GetByName(ctx, userName)
	if !ok {
		return nil, false, ErrTechnical
	}
	if user == nil {
		return nil, false, ErrProfileNotFound
	}

	if requestingUserName == "" {
		return user, false, nil
	}

	reqUser, ok := i.userRW.GetByName(ctx, requestingUserName)
	if !ok {
		return nil, false, ErrTechnical
	}
	if user == nil {
		return nil, false, ErrProfileNotFound
	}

	return user, reqUser.Follows(userName), nil
}
