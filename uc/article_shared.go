package uc

import (
	"errors"

	"github.com/err0r500/go-realworld-clean/domain"
)

func (i interactor) getArticleAndCheckUser(username, slug string) (*domain.Article, error) {
	completeArticle, err := i.articleRW.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if completeArticle == nil {
		return nil, errArticleNotFound
	}

	// check only if a username is specified
	if username != "" && completeArticle.Author.Name != username {
		return nil, errors.New("article not owned by user")
	}

	return completeArticle, nil
}
