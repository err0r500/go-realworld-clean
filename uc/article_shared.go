package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func (i interactor) getArticleAndCheckUser(ctx context.Context, username, slug string) (*domain.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:get_article_check_user")
	defer span.Finish()

	completeArticle, ok := i.articleRW.GetBySlug(ctx, slug)
	if !ok {
		return nil, errTechnical
	}
	if completeArticle == nil {
		return nil, errArticleNotFound
	}

	// check only if a username is specified
	if username != "" && completeArticle.Author.Name != username {
		span.LogFields(log.Error(ErrUnauthorized))
		return nil, ErrUnauthorized
	}

	return completeArticle, nil
}
