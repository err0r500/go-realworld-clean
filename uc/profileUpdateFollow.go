package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) ProfileUpdateFollow(ctx context.Context, userName, followeeName string, follow bool) (*domain.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:profile_update_follow")
	defer span.Finish()

	mayUser, ok := i.userRW.GetByName(ctx, userName)
	if !ok {
		return nil, ErrTechnical
	}
	if mayUser == nil {
		return nil, ErrNotFound
	}

	if mayUser.Name != userName {
		return nil, ErrUnauthorized
	}

	mayUser.UpdateFollowees(followeeName, follow)

	if ok := i.userRW.Save(ctx, *mayUser); !ok {
		return nil, ErrTechnical
	}

	return mayUser, nil
}
