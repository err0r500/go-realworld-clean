package uc

import (
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

func (i interactor) GetArticles(limit, offset int, filters []domain.ArticleFilter) (domain.ArticleCollection, int, error) {
	if limit <= 0 {
		return domain.ArticleCollection{}, 0, nil
	}

	articles, err := i.articleRW.GetRecentFiltered(filters)
	if err != nil {
		return nil, 0, err
	}

	return domain.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset), len(articles), nil
}
