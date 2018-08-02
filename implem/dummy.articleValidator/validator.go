package articleValidator

import (
	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/uc"
)

type validator struct {
}

func New() uc.ArticleValidator {
	return validator{}
}

func (validator) BeforeCreationCheck(article *domain.Article) error { return nil }
func (validator) BeforeUpdateCheck(article *domain.Article) error   { return nil }
