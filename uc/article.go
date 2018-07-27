package uc

import "github.com/err0r500/go-realworld-clean/domain"

func (i interactor) ArticlePost(username string, article domain.Article) (*domain.Article, error) {
	return nil, nil
}

func (i interactor) ArticlePut(username string, slug string, article domain.Article) (*domain.Article, error) {
	return nil, nil
}

func (i interactor) ArticleGet(slug string) (*domain.Article, error) {
	return nil, nil
}

func (i interactor) ArticleDelete(username string, slug string) error {
	return nil
}
