package uc_test

import (
	"testing"

	"errors"

	"github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInteractor_CommentsGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	article := testData.Article("jane")

	i := mock.NewMockedInteractor(mockCtrl)
	i.ArticleRW.EXPECT().GetBySlug(article.Slug).Return(&article, nil)

	comments, err := i.GetUCHandler().CommentsGet(article.Slug)
	assert.NoError(t, err)
	assert.Equal(t, article.Comments, comments)
}

func TestInteractor_CommentsGet_fails(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	article := testData.Article("jane")

	i := mock.NewMockedInteractor(mockCtrl)
	i.ArticleRW.EXPECT().GetBySlug(article.Slug).Return(nil, errors.New(""))

	comments, err := i.GetUCHandler().CommentsGet(article.Slug)
	assert.Error(t, err)
	assert.Nil(t, comments)
}
