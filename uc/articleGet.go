package uc

import "github.com/err0r500/go-realworld-clean/domain"

func (i interactor) ArticleGet(username, slug string) (*domain.User, *domain.Article, error) {
	var user *domain.User
	if username != "" {
		var errGet error
		user, errGet = i.userRW.GetByName(username)
		if errGet != nil {
			return nil, nil, errGet
		}
	}

	article, err := i.articleRW.GetBySlug(slug)
	if err != nil {
		return nil, nil, err
	}

	return user, article, nil
}
