package articleRW

import (
	"sync"

	"errors"

	"time"

	"log"

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
func (rw rw) Create(article domain.Article) (*domain.Article, error) {
	if _, err := rw.GetBySlug(article.Slug); err == nil {
		log.Println(err)
		return nil, uc.ErrAlreadyInUse
	}
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	rw.store.Store(article.Slug, article)

	return rw.GetBySlug(article.Slug)
}

func (rw rw) GetByAuthorsNameOrderedByMostRecentAsc(usernames []string) ([]domain.Article, error) {
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

	return toReturn, nil
}

func (rw rw) GetRecentFiltered(filters []domain.ArticleFilter) ([]domain.Article, error) {
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

	return recentArticles, nil
}

func (rw rw) Save(article domain.Article) (*domain.Article, error) {
	if _, err := rw.GetBySlug(article.Slug); err != nil {
		return nil, uc.ErrNotFound
	}

	article.UpdatedAt = time.Now()
	rw.store.Store(article.Slug, article)

	return rw.GetBySlug(article.Slug)
}

func (rw rw) GetBySlug(slug string) (*domain.Article, error) {
	value, ok := rw.store.Load(slug)
	if !ok {
		return nil, uc.ErrNotFound
	}

	article, ok := value.(domain.Article)
	if !ok {
		return nil, errors.New("not an article stored at key")
	}

	return &article, nil
}

func (rw rw) Delete(slug string) error {
	rw.store.Delete(slug)

	return nil
}
