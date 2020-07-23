package uc

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/err0r500/go-realworld-clean/domain"
)

func NewFilters(author, tag, favorite string) []domain.ArticleFilter {
	var filters []domain.ArticleFilter
	if author != "" {
		filters = append(filters, domain.ArticleHasAuthor(author))
	}
	if tag != "" {
		filters = append(filters, domain.ArticleHasTag(tag))
	}
	if favorite != "" {
		filters = append(filters, domain.ArticleIsFavoritedBy(favorite))
	}

	return filters
}

func (i interactor) GetArticles(ctx context.Context, username string, limit, offset int, filters []domain.ArticleFilter) (*domain.User, domain.ArticleCollection, int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:articles_get")
	defer span.Finish()

	if limit <= 0 {
		return nil, domain.ArticleCollection{}, 0, nil
	}

	articles, err := i.articleRW.GetRecentFiltered(ctx, filters)
	if err != nil {
		return nil, nil, 0, err
	}

	var user *domain.User
	if username != "" {
		var errGet error
		user, errGet = i.userRW.GetByName(username)
		if errGet != nil {
			return nil, nil, 0, errGet
		}
	}

	return user, domain.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset), len(articles), nil
}
