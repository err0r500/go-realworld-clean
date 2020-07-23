package uc

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

func (i interactor) Tags(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:tags")
	defer span.Finish()

	tags, ok := i.tagsRW.GetAll(ctx)
	if !ok {
		return nil, ErrTechnical
	}
	return tags, nil
}
