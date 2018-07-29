package uc

import (
	"github.com/err0r500/go-realworld-clean/domain"
)

func (i interactor) CommentsGet(slug string) ([]domain.Comment, error) {
	article, err := i.articleRW.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	return article.Comments, nil
}

func (i interactor) CommentsPost(username, slug, comment string) (*domain.Comment, error) {
	commentPoster, err := i.userRW.GetByName(username)
	if err != nil {
		return nil, err
	}

	article, err := i.articleRW.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	rawComment := domain.Comment{
		Body:   comment,
		Author: *commentPoster,
	}

	insertedComment, err := i.commentRW.Create(rawComment)
	if err != nil {
		return nil, err
	}

	article.Comments = append(article.Comments, *insertedComment)

	if _, err := i.articleRW.Save(*article); err != nil {
		return nil, err
	}

	return insertedComment, nil
}

func (i interactor) CommentsDelete(username, slug string, id int) error {
	comment, err := i.commentRW.GetByID(id)
	if err != nil {
		return err
	}
	if comment.Author.Name != username {
		return errWrongUser
	}

	if err := i.commentRW.Delete(id); err != nil {
		return err
	}

	article, err := i.articleRW.GetBySlug(slug)
	if err != nil {
		return err
	}

	article.UpdateComments(*comment, false)

	if _, err := i.articleRW.Save(*article); err != nil {
		return err
	}

	return nil
}
