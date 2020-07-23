package uc

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

func (i interactor) ArticleDelete(ctx context.Context, username string, slug string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:article_delete")
	defer span.Finish()

	if _, err := i.getArticleAndCheckUser(ctx, username, slug); err != nil {
		return err
	}

	if ok := i.articleRW.Delete(ctx, slug); !ok {
		return errTechnical
	}

	return nil
}
