package uc

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/err0r500/go-realworld-clean/domain"
)

func (i interactor) ArticlePut(ctx context.Context, username string, slug string, fieldsToUpdate map[domain.ArticleUpdatableField]*string) (*domain.User, *domain.Article, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "uc:article_put")
	defer span.Finish()

	article, err := i.getArticleAndCheckUser(ctx, username, slug)
	if err != nil {
		return nil, nil, err
	}

	domain.UpdateArticle(article,
		domain.SetArticleTitle(fieldsToUpdate[domain.Title]),
		domain.SetArticleDescription(fieldsToUpdate[domain.Description]),
		domain.SetArticleBody(fieldsToUpdate[domain.Body]),
	)

	if err := i.articleValidator.BeforeUpdateCheck(article); err != nil {
		return nil, nil, err
	}

	user, ok := i.userRW.GetByName(ctx, username)
	if !ok {
		return nil, nil, ErrTechnical
	}
	if user == nil {
		return nil, nil, ErrNotFound
	}

	savedArticle, ok := i.articleRW.Save(ctx, *article)
	if !ok {
		return nil, nil, ErrTechnical
	}

	return user, savedArticle, nil
}
