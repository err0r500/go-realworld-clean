package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) ArticlePost(ctx context.Context, username string, article domain.Article) (*domain.User, *domain.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:article_post")
	defer span.Finish()

	user, ok := i.userRW.GetByName(ctx, username)
	if !ok {
		return nil, nil, ErrTechnical
	}
	if user == nil {
		return nil, nil, ErrUnauthorized
	}

	slug := i.slugger.NewSlug(article.Title)
	art, ok := i.articleRW.GetBySlug(ctx, slug)
	if !ok {
		return nil, nil, ErrTechnical
	}
	if art != nil {
		return nil, nil, ErrConflict
	}

	article.Slug = slug
	article.Author = *user

	if err := i.articleValidator.BeforeCreationCheck(&article); err != nil {
		return nil, nil, ErrValidation
	}

	completeArticle, ok := i.articleRW.Create(ctx, article)
	if !ok {
		return nil, nil, ErrTechnical
	}

	if ok := i.tagsRW.Add(ctx, article.TagList); !ok {
		return nil, nil, ErrTechnical
	}

	return user, completeArticle, nil
}
