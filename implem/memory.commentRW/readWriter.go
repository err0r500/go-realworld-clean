package commentRW

import (
	"context"
	"sync"

	"github.com/opentracing/opentracing-go/log"

	"github.com/opentracing/opentracing-go"

	"time"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/uc"
)

type rw struct {
	store *sync.Map
}

func New() uc.CommentRW {
	return rw{
		store: &sync.Map{},
	}
}

func (rw rw) Create(ctx context.Context, comment domain.Comment) (*domain.Comment, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentrw:create")
	defer span.Finish()

	found, ok := rw.GetByID(ctx, comment.ID)
	if !ok {
		span.LogFields(log.Error(uc.ErrTechnical))
		return nil, false
	}
	if found != nil {
		span.LogFields(log.Error(uc.ErrConflict))
		return nil, false
	}

	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	rw.store.Store(comment.ID, comment)

	return &comment, true
}

func (rw rw) GetByID(ctx context.Context, id int) (*domain.Comment, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentrw:get_by_id")
	defer span.Finish()

	value, ok := rw.store.Load(id)
	if !ok {
		span.LogFields(log.Error(uc.ErrNotFound))
		return nil, false
	}

	comment, ok := value.(domain.Comment)
	if !ok {
		span.LogFields(log.Error(uc.ErrTechnical))
		return nil, false
	}

	return &comment, true
}

func (rw rw) Delete(ctx context.Context, id int) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commentrw:get_by_id")
	defer span.Finish()

	rw.store.Delete(id)

	return true
}
