package uc

import (
	"strconv"

	"github.com/err0r500/go-realworld-clean/domain"
)

type Filters struct {
	AuthorFilter    *string
	TagFilter       *string
	FavoritedFilter *bool
}

func NewFilters(author, tag, favorite string) Filters {
	filters := Filters{}
	if author != "" {
		filters.AuthorFilter = &author
	}
	if tag != "" {
		filters.TagFilter = &tag
	}

	fav, err := strconv.ParseBool(favorite)
	if err != nil {
		return filters
	}
	filters.FavoritedFilter = &fav

	return filters
}

func (i interactor) GetArticles(limit, offset int, filters Filters) (domain.ArticleCollection, int, error) {
	if limit <= 0 {
		return domain.ArticleCollection{}, 0, nil
	}

	articles, err := i.articleRW.GetRecentFiltered(filters)
	if err != nil {
		return nil, 0, err
	}

	return domain.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset), len(articles), nil
}
