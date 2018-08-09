package commentRW

import (
	"sync"

	"errors"

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

func (rw rw) Create(comment domain.Comment) (*domain.Comment, error) {
	if _, err := rw.GetByID(comment.ID); err == nil {
		return nil, uc.ErrAlreadyInUse
	}
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	rw.store.Store(comment.ID, comment)

	return rw.GetByID(comment.ID)
}

func (rw rw) GetByID(id int) (*domain.Comment, error) {
	value, ok := rw.store.Load(id)
	if !ok {
		return nil, uc.ErrNotFound
	}

	comment, ok := value.(domain.Comment)
	if !ok {
		return nil, errors.New("not an article stored at key")
	}

	return &comment, nil
}

func (rw rw) Delete(id int) error {
	rw.store.Delete(id)

	return nil
}
