package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) CommentsPost(ctx context.Context, username, slug, comment string) (*domain.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:comments_post")
	defer span.Finish()

	commentPoster, err := i.userRW.GetByName(username)
	if err != nil {
		return nil, err
	}

	article, ok := i.articleRW.GetBySlug(ctx, slug)
	if !ok {
		return nil, errTechnical
	}
	if article == nil {
		return nil, ErrNotFound
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

	if _, ok := i.articleRW.Save(ctx, *article); !ok {
		return nil, ErrTechnical
	}

	return insertedComment, nil
}
