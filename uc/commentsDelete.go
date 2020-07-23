package uc

import (
	"context"

	"github.com/opentracing/opentracing-go/log"

	"github.com/opentracing/opentracing-go"
)

func (i interactor) CommentsDelete(ctx context.Context, username, slug string, id int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:comments_delete")
	defer span.Finish()

	comment, ok := i.commentRW.GetByID(ctx, id)
	if !ok {
		return ErrTechnical
	}
	if comment.Author.Name != username {
		return ErrUnauthorized
	}

	if ok := i.commentRW.Delete(ctx, id); !ok {
		return ErrTechnical
	}

	article, ok := i.articleRW.GetBySlug(ctx, slug)
	if !ok {
		return ErrTechnical
	}
	if article == nil {
		span.LogFields(log.Error(ErrNotFound))
		return ErrNotFound
	}

	article.UpdateComments(*comment, false)

	if _, ok := i.articleRW.Save(ctx, *article); !ok {
		return ErrTechnical
	}

	return nil
}
