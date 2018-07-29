package uc

import "github.com/err0r500/go-realworld-clean/domain"

func (i interactor) FavoritesUpdate(username, slug string, favorite bool) (*domain.Article, error) {
	article, err := i.articleRW.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	user, err := i.userRW.GetByName(username)
	if err != nil {
		return nil, err
	}

	user.UpdateFavorites(*article, favorite)
	if err := i.userRW.Save(*user); err != nil {
		return nil, err
	}

	return article, nil
}
