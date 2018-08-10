package uc

import (
	"github.com/err0r500/go-realworld-clean/domain"
)

func (i interactor) ArticlePost(username string, article domain.Article) (*domain.User, *domain.Article, error) {
	user, err := i.userRW.GetByName(username)
	if err != nil {
		return nil, nil, err
	}

	slug := i.slugger.NewSlug(article.Title)
	if _, err := i.articleRW.GetBySlug(slug); err == nil {
		return nil, nil, ErrAlreadyInUse
	}

	article.Slug = slug
	article.Author = *user

	if err := i.articleValidator.BeforeCreationCheck(&article); err != nil {
		return nil, nil, err
	}

	completeArticle, err := i.articleRW.Create(article)
	if err != nil {
		return nil, nil, err
	}

	if err := i.tagsRW.Add(article.TagList); err != nil {
		return nil, nil, err
	}

	return user, completeArticle, nil
}
