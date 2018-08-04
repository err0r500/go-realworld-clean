package uc

import (
	"github.com/err0r500/go-realworld-clean/domain"
)

func NewFilters(author, tag, favorite string) []domain.ArticleFilter {
	filters := []domain.ArticleFilter{}
	if author != "" {
		filters = append(filters, domain.HasAuthor(author))
	}
	if tag != "" {
		filters = append(filters, domain.Hastag(tag))
	}

	//fav, err := strconv.ParseBool(favorite)
	//if err != nil {
	//	return filters
	//}
	//filters.FavoritedFilter = &fav

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
