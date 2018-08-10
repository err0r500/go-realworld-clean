package uc_test

import (
	"testing"

	"errors"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInteractor_ArticleDelete_happycase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	jane := testData.User("jane")
	article := testData.Article("janeArticle")

	i := mock.NewMockedInteractor(mockCtrl)
	i.ArticleRW.EXPECT().GetBySlug(article.Slug).Return(&article, nil).Times(1)
	i.ArticleRW.EXPECT().Delete(article.Slug)

	assert.NoError(t, i.GetUCHandler().ArticleDelete(jane.Name, article.Slug))
}

func TestInteractor_ArticleDelete_fails(t *testing.T) {
	mutations := map[string]mock.Tester{
		"shouldPass": {
			Calls: func(i *mock.Interactor) { // change nothing
			},
			ShouldPass: true},
		"error return on aRW.GetBySlug": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(nil, errors.New(""))
			}},
		"wrong author return on aRW.GetBySlug": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(&domain.Article{Author: domain.User{Name: "heyhey"}}, nil)
			}},
		"error return on delete": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().Delete(gomock.Any()).Return(errors.New(""))
			}},
	}
	jane := testData.User("jane")
	article := testData.Article("janeArticle")

	// same as the happy case but with any parameter and called any number of times (including 0)
	validCalls := func(i *mock.Interactor) {
		i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(&article, nil).AnyTimes()
		i.ArticleRW.EXPECT().Delete(gomock.Any()).AnyTimes()
	}

	for testName, mutation := range mutations {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			i := mock.NewMockedInteractor(mockCtrl)
			mutation.Calls(&i) // put the tested call first (important)
			validCalls(&i)     // then fill the gaps with valid calls

			err := i.GetUCHandler().ArticleDelete(jane.Name, article.Slug)
			if mutation.ShouldPass {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}
