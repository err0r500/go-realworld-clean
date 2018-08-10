package uc_test

import (
	"testing"

	"errors"

	"github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInteractor_ArticleGet_happycase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	testUser := testData.User("rick")
	article := testData.Article("janeArticle")

	i := mock.NewMockedInteractor(mockCtrl)
	i.ArticleRW.EXPECT().GetBySlug(article.Slug).Return(&article, nil).Times(1)
	i.UserRW.EXPECT().GetByName(testUser.Name).Return(&testUser, nil).Times(1)

	_, _, err := i.GetUCHandler().ArticleGet(testUser.Name, article.Slug)
	assert.NoError(t, err)
}

func TestInteractor_ArticleGet_fails(t *testing.T) {
	mutations := map[string]mock.Tester{
		"shouldPass": {
			Calls: func(i *mock.Interactor) { // change nothing
			},
			ShouldPass: true},
		"error return on uRW.GetByName": {
			Calls: func(i *mock.Interactor) {
				i.UserRW.EXPECT().GetByName(gomock.Any()).Return(nil, errors.New(""))
			}},
		"error return on aRW.GetBySlug": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(nil, errors.New(""))
			}},
	}

	testUser := testData.User("rick")
	article := testData.Article("janeArticle")

	// same as the happy case but with any parameter and called any number of times (including 0)
	validCalls := func(i *mock.Interactor) {
		i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(&article, nil).AnyTimes()
		i.UserRW.EXPECT().GetByName(gomock.Any()).Return(&testUser, nil).AnyTimes()
	}

	for testName, mutation := range mutations {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			i := mock.NewMockedInteractor(mockCtrl)
			mutation.Calls(&i) // put the tested call first (important)
			validCalls(&i)     // then fill the gaps with valid calls

			_, _, err := i.GetUCHandler().ArticleGet(testUser.Name, article.Slug)
			if mutation.ShouldPass {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}
