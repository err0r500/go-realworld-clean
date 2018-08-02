package tagsRW

import (
	"sync"

	"github.com/err0r500/go-realworld-clean/uc"
)

type rw struct {
	store *sync.Map
}

func New() uc.TagsRW {
	return rw{
		store: &sync.Map{},
	}
}

// lots of ways to improve this (use an array as cache, use index access instead of append...)
// no perf problem for now => no optimisation :)
func (rw rw) GetAll() ([]string, error) {
	var toReturn []string

	rw.store.Range(func(key, value interface{}) bool {
		tag, ok := key.(string)
		if !ok {
			return true
		}
		toReturn = append(toReturn, tag)
		return true
	})

	return toReturn, nil
}

func (rw rw) Add(newTags []string) error {

	for _, tag := range newTags {
		rw.store.Store(tag, true)
	}

	return nil
}
