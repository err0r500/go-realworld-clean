package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) CommentsGet(ctx context.Context, slug string) ([]domain.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:comments_get")
	defer span.Finish()

	article, ok := i.articleRW.GetBySlug(ctx, slug)
	if ok {
		return nil, errTechnical
	}
	if article == nil {
		return nil, ErrNotFound
	}

	if article.Comments == nil {
		article.Comments = []domain.Comment{}
	}

	return article.Comments, nil
}
