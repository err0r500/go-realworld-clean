package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/opentracing/opentracing-go"
)

func (i interactor) ArticlePost(ctx context.Context, username string, article domain.Article) (*domain.User, *domain.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:article_post")
	defer span.Finish()

	user, err := i.userRW.GetByName(username)
	if err != nil {
		return nil, nil, err
	}

	slug := i.slugger.NewSlug(article.Title)
	art, ok := i.articleRW.GetBySlug(ctx, slug)
	if !ok {
		return nil, nil, errTechnical
	}
	if art != nil {
		return nil, nil, ErrConflict
	}

	article.Slug = slug
	article.Author = *user

	if err := i.articleValidator.BeforeCreationCheck(&article); err != nil {
		return nil, nil, err
	}

	completeArticle, ok := i.articleRW.Create(ctx, article)
	if !ok {
		return nil, nil, ErrTechnical
	}

	if err := i.tagsRW.Add(article.TagList); err != nil {
		return nil, nil, err
	}

	return user, completeArticle, nil
}
