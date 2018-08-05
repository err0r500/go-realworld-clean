package uc

import "github.com/err0r500/go-realworld-clean/domain"

// fixme : should return total favorite count and fav for current user
func (i interactor) FavoritesUpdate(username, slug string, favorite bool) (*domain.Article, error) {
	article, err := i.articleRW.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	user, err := i.userRW.GetByName(username)
	if err != nil {
		return nil, err
	}

	article.UpdateFavoritedBy(*user, favorite)

	if err := i.userRW.Save(*user); err != nil {
		return nil, err
	}

	return article, nil
}
