package uc

import "github.com/err0r500/go-realworld-clean/domain"

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
