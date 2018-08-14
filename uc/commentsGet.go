package uc

import "github.com/err0r500/go-realworld-clean/domain"

func (i interactor) CommentsGet(slug string) ([]domain.Comment, error) {
	article, err := i.articleRW.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if article.Comments == nil {
		article.Comments = []domain.Comment{}
	}

	return article.Comments, nil
}
