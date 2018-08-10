package uc

import "github.com/err0r500/go-realworld-clean/domain"

func (i interactor) FavoritesUpdate(username, slug string, favorite bool) (*domain.User, *domain.Article, error) {
	user, err := i.userRW.GetByName(username)
	if err != nil {
		return nil, nil, err
	}

	article, err := i.articleRW.GetBySlug(slug)
	if err != nil {
		return nil, nil, err
	}

	article.UpdateFavoritedBy(*user, favorite)

	updatedArticle, err := i.articleRW.Save(*article)
	if err != nil {
		return nil, nil, err
	}

	return user, updatedArticle, nil
}
