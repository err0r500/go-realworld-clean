package uc_test

import (
	"testing"

	"errors"

	"github.com/err0r500/go-realworld-clean/domain"
	mock "github.com/err0r500/go-realworld-clean/implem/mock.uc"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInteractor_ArticlePost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	rick := testData.User("rick")

	article := domain.Article{
		Title: "article Title",
	}
	slug := "article-title"

	i := mock.NewMockedInteractor(mockCtrl)
	i.UserRW.EXPECT().GetByName(rick.Name).Return(&rick, nil).Times(1)
	i.Slugger.EXPECT().NewSlug(article.Title).Return(slug).Times(1)
	i.ArticleRW.EXPECT().GetBySlug(slug).Return(nil, errors.New("not found")).Times(1)
	i.ArticleValidator.EXPECT().BeforeCreationCheck(gomock.Any()).Return(nil).Times(1)
	i.ArticleRW.EXPECT().Create(gomock.Any())

	_, err := i.GetUCHandler().ArticlePost(rick.Name, article)
	assert.NoError(t, err)
}

func TestInteractor_ArticlePut(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	jane := testData.User("jane")
	origArticle := testData.Article("janeArticle")
	req := domain.Article{
		Title:       "newTitle",
		Description: "newDescription",
		Body:        origArticle.Body,
		TagList:     origArticle.TagList,
	}

	toInsert := origArticle
	toInsert.Title = req.Title
	toInsert.Description = req.Description

	i := mock.NewMockedInteractor(mockCtrl)
	i.ArticleRW.EXPECT().GetBySlug(origArticle.Slug).Return(&origArticle, nil).Times(1)
	i.ArticleValidator.EXPECT().BeforeUpdateCheck(&toInsert).Return(nil).Times(1)
	i.ArticleRW.EXPECT().Save(toInsert)

	_, err := i.GetUCHandler().ArticlePut(jane.Name, origArticle.Slug, req)
	assert.NoError(t, err)
}

func TestInteractor_ArticleDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	jane := testData.User("jane")
	article := testData.Article("janeArticle")

	i := mock.NewMockedInteractor(mockCtrl)
	i.ArticleRW.EXPECT().GetBySlug(article.Slug).Return(&article, nil).Times(1)
	i.ArticleRW.EXPECT().Delete(article.Slug)

	assert.NoError(t, i.GetUCHandler().ArticleDelete(jane.Name, article.Slug))
}

func TestInteractor_ArticleGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	article := testData.Article("janeArticle")

	i := mock.NewMockedInteractor(mockCtrl)
	i.ArticleRW.EXPECT().GetBySlug(article.Slug).Return(&article, nil).Times(1)

	_, err := i.GetUCHandler().ArticleGet(article.Slug)
	assert.NoError(t, err)
}
