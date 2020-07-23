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

	articles, ok := i.articleRW.GetRecentFiltered(ctx, filters)
	if !ok {
		return nil, nil, 0, ErrTechnical
	}

	var user *domain.User
	if username != "" {
		user, ok = i.userRW.GetByName(ctx, username)
		if !ok {
			return nil, nil, 0, ErrTechnical
		}
	}

	return user, domain.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset), len(articles), nil
}
