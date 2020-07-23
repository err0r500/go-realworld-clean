package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) ProfileUpdateFollow(ctx context.Context, userName, followeeName string, follow bool) (*domain.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:profile_update_follow")
	defer span.Finish()

	user, ok := i.userRW.GetByName(ctx, userName)
	if !ok {
		return nil, ErrTechnical
	}
	if user.Name != userName {
		return nil, ErrUnauthorized
	}
	if user == nil {
		return nil, ErrNotFound
	}

	user.UpdateFollowees(followeeName, follow)

	if ok := i.userRW.Save(ctx, *user); !ok {
		return nil, ErrTechnical
	}

	return user, nil
}
