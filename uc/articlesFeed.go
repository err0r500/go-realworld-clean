package uc

import (
	"github.com/err0r500/go-realworld-clean/domain"
)

func (i interactor) ArticlesFeed(username string, limit, offset int) (domain.ArticleCollection, int, error) {
	if limit < 0 {
		return domain.ArticleCollection{}, 0, nil
	}

	user, err := i.userRW.GetByName(username)
	if err != nil {
		return nil, 0, err
	}
	articles, err := i.articleRW.GetByAuthorsNameOrderedByMostRecentAsc(user.FollowIDs)
	if err != nil {
		return nil, 0, err
	}

	return domain.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset), len(articles), nil
}
