package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) FavoritesUpdate(ctx context.Context, username, slug string, favorite bool) (*domain.User, *domain.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:favorites_update")
	defer span.Finish()

	user, err := i.userRW.GetByName(username)
	if err != nil {
		return nil, nil, err
	}

	article, ok := i.articleRW.GetBySlug(ctx, slug)
	if !ok {
		return nil, nil, errTechnical
	}
	if article == nil {
		return nil, nil, ErrNotFound
	}

	article.UpdateFavoritedBy(*user, favorite)

	updatedArticle, ok := i.articleRW.Save(ctx, *article)
	if !ok {
		return nil, nil, ErrTechnical
	}

	return user, updatedArticle, nil
}
