package uc

import (
	"errors"

	"github.com/err0r500/go-realworld-clean/domain"
)

func (i interactor) ArticlePost(username string, article domain.Article) (*domain.Article, error) {
	user, err := i.userRW.GetByName(username)
	if err != nil {
		return nil, err
	}

	slug := i.slugger.NewSlug(article.Title)
	if _, err := i.articleRW.GetBySlug(slug); err == nil {
		return nil, ErrAlreadyInUse
	}

	article.Slug = slug
	article.Author = *user

	if err := i.articleValidator.BeforeCreationCheck(&article); err != nil {
		return nil, err
	}

	completeArticle, err := i.articleRW.Create(article)
	if err != nil {
		return nil, err
	}

	if err := i.tagsRW.Add(article.TagList); err != nil {
		return nil, err
	}

	return completeArticle, nil
}

func (i interactor) ArticlePut(username string, slug string, reqArticle domain.Article) (*domain.Article, error) {
	article, err := i.getArticleAndCheckUser(username, slug)
	if err != nil {
		return nil, err
	}

	// real PUT request, all fields are mandatory in request
	article.Title = reqArticle.Title
	article.Description = reqArticle.Description
	article.Body = reqArticle.Body
	article.TagList = reqArticle.TagList

	if err := i.articleValidator.BeforeUpdateCheck(article); err != nil {
		return nil, err
	}

	return i.articleRW.Save(*article)
}

func (i interactor) ArticleGet(slug string) (*domain.Article, error) {
	return i.getArticleAndCheckUser("", slug)
}

func (i interactor) ArticleDelete(username string, slug string) error {
	_, err := i.getArticleAndCheckUser(username, slug)
	if err != nil {
		return err
	}

	return i.articleRW.Delete(slug)
}

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
