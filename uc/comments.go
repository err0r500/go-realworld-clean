package uc

import "github.com/err0r500/go-realworld-clean/domain"

func (i interactor) CommentsGet(slug string) ([]domain.Comment, error) {
	return nil, nil
}

func (i interactor) CommentsPost(username, slug, comment string) (*domain.Comment, error) {
	return nil, nil
}

func (i interactor) CommentsDelete(username string, slug, id string) error {
	return nil
}
