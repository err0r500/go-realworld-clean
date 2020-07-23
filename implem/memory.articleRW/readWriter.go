package articleRW

import (
	"context"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"errors"

	"time"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/uc"
)

type rw struct {
	store *sync.Map
}

func New() uc.ArticleRW {
	return rw{
		store: &sync.Map{},
	}
}
func (rw rw) Create(ctx context.Context, article domain.Article) (*domain.Article, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inmem_articlerw:create")
	defer span.Finish()

	art, ok := rw.GetBySlug(ctx, article.Slug)
	if !ok {
		span.LogFields(log.Error(uc.ErrTechnical))
		return nil, false
	}
	if art != nil {
		span.LogFields(log.Error(uc.ErrConflict))
		return nil, false
	}

	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	rw.store.Store(article.Slug, article)

	return rw.GetBySlug(ctx, article.Slug)
}

func (rw rw) GetByAuthorsNameOrderedByMostRecentAsc(ctx context.Context, usernames []string) ([]domain.Article, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inmem_articlerw:get_by_author_most_recent")
	defer span.Finish()

	var toReturn []domain.Article

	rw.store.Range(func(key, value interface{}) bool {
		article, ok := value.(domain.Article)
		if !ok {
			return true // log this but continue
		}
		for _, username := range usernames {
			if article.Author.Name == username {
				toReturn = append(toReturn, article)
			}
		}
		return true
	})

	return toReturn, true
}

func (rw rw) GetRecentFiltered(ctx context.Context, filters []domain.ArticleFilter) ([]domain.Article, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inmem_articlerw:get_recent_filtered")
	defer span.Finish()

	var recentArticles []domain.Article

	rw.store.Range(func(key, value interface{}) bool {
		article, ok := value.(domain.Article)
		if !ok {
			// not an article (shouldn't happen) -> skip
			return true // log this but continue
		}

		for _, funcToApply := range filters {
			if !funcToApply(article) { // "AND filter" : if one of the filter is at false, skip the article
				return true
			}
		}

		recentArticles = append(recentArticles, article)
		return true
	})

	return recentArticles, true
}

func (rw rw) Save(ctx context.Context, article domain.Article) (*domain.Article, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inmem_articlerw:save")
	defer span.Finish()

	art, ok := rw.GetBySlug(ctx, article.Slug)
	if !ok {
		span.LogFields(log.Error(uc.ErrTechnical))
		return nil, false
	}
	if art == nil {
		span.LogFields(log.Error(uc.ErrNotFound))
		return nil, false
	}

	article.UpdatedAt = time.Now()
	rw.store.Store(article.Slug, article)

	return rw.GetBySlug(ctx, article.Slug)
}

func (rw rw) GetBySlug(ctx context.Context, slug string) (*domain.Article, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inmem_articlerw:get_by_slug")
	defer span.Finish()

	value, ok := rw.store.Load(slug)
	if !ok {
		return nil, true
	}

	article, ok := value.(domain.Article)
	if !ok {
		span.LogFields(log.Error(errors.New("not an article stored at key")))
		return nil, false
	}

	return &article, true
}

func (rw rw) Delete(ctx context.Context, slug string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "inmem_articlerw:delete")
	defer span.Finish()

	rw.store.Delete(slug)

	return true
}
