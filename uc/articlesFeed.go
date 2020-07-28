package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) ArticlesFeed(ctx context.Context, username string, limit, offset int) (*domain.User, domain.ArticleCollection, int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:article_get_feed")
	defer span.Finish()

	if limit < 0 {
		return nil, domain.ArticleCollection{}, 0, nil
	}

	var user *domain.User
	if username != "" {
		mayUser, ok := i.userRW.GetByName(ctx, username)
		if !ok {
			return nil, nil, 0, ErrTechnical
		}
		if mayUser == nil {
			return nil, nil, 0, ErrNotFound
		}
		user = mayUser
	}

	articles, ok := i.articleRW.GetByAuthorsNameOrderedByMostRecentAsc(ctx, user.FollowIDs)
	if !ok {
		return nil, nil, 0, ErrTechnical
	}

	return user, domain.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset), len(articles), nil // needs the original length
}
