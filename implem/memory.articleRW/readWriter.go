package articleRW

import (
	"sync"

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

func (rw) GetByAuthorsNameOrderedByMostRecentAsc(usernames []string) ([]domain.Article, error) {
	return nil, nil
}

func (rw) GetRecentFiltered(filters uc.Filters) ([]domain.Article, error) { return nil, nil }

func (rw) Create(domain.Article) (*domain.Article, error) { return nil, nil }
func (rw) Save(domain.Article) (*domain.Article, error)   { return nil, nil }
func (rw) GetBySlug(slug string) (*domain.Article, error) { return nil, nil }
func (rw) Delete(slug string) error                       { return nil }
