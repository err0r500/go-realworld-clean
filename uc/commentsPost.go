package uc

import (
	"context"
	"errors"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func (i interactor) CommentsPost(ctx context.Context, username, slug, comment string) (*domain.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:comments_post")
	defer span.Finish()

	commentPoster, ok := i.userRW.GetByName(ctx, username)
	if !ok {
		return nil, ErrTechnical
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

	insertedComment, ok := i.commentRW.Create(ctx, rawComment)
	if !ok {
		return nil, ErrTechnical
	}
	if insertedComment == nil {
		span.LogFields(log.Error(errors.New("create comment returned a nil pointer")))
		return nil, ErrTechnical
	}

	article.Comments = append(article.Comments, *insertedComment)

	if _, ok := i.articleRW.Save(ctx, *article); !ok {
		return nil, ErrTechnical
	}

	return insertedComment, nil
}
