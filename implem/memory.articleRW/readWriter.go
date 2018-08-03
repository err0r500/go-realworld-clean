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

func (rw) GetRecentFiltered(filters uc.Filters) ([]domain.Article, error) {
	// todo => check if its AND or OR filters

	return nil, nil
}

func (rw rw) Save(article domain.Article) (*domain.Article, error) {
	if _, err := rw.GetBySlug(article.Slug); err != nil {
		return nil, uc.ErrNotFound
	}

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
