package uc

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/err0r500/go-realworld-clean/domain"
)

func (i interactor) ArticleGet(ctx context.Context, username, slug string) (*domain.User, *domain.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:article_get")
	defer span.Finish()

	var user *domain.User
	if username != "" {
		mayUser, ok := i.userRW.GetByName(ctx, username)
		if !ok {
			return nil, nil, ErrTechnical
		}
		if mayUser == nil {
			return nil, nil, ErrNotFound
		}
		user = mayUser
	}

	article, ok := i.articleRW.GetBySlug(ctx, slug)
	if !ok {
		return nil, nil, errTechnical
	}

	return user, article, nil
}
