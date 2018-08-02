package slugger

import (
	"github.com/err0r500/go-realworld-clean/uc"
	"github.com/gosimple/slug"
)

type slugger struct{}

func New() uc.Slugger {
	return slugger{}
}

func (slugger) NewSlug(initial string) string {
	return slug.Make(initial)
}
