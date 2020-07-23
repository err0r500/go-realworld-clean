package uc

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

func (i interactor) CommentsDelete(ctx context.Context, username, slug string, id int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:comments_delete")
	defer span.Finish()

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

	article, ok := i.articleRW.GetBySlug(ctx, slug)
	if !ok {
		return err
	}

	article.UpdateComments(*comment, false)

	if _, ok := i.articleRW.Save(ctx, *article); !ok {
		return ErrTechnical
	}

	return nil
}
