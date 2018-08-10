package uc_test

import (
	"errors"
	"testing"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInteractor_ArticlePost_happycase(t *testing.T) {
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
	i.TagsRW.EXPECT().Add(gomock.Any())
	i.ArticleRW.EXPECT().Create(gomock.Any())

	_, _, err := i.GetUCHandler().ArticlePost(rick.Name, article)
	assert.NoError(t, err)
}

func TestInteractor_ArticlePost_fails(t *testing.T) {
	mutations := map[string]mock.Tester{
		"shouldPass": {
			Calls: func(i *mock.Interactor) { // change nothing
			},
			ShouldPass: true},
		"error return on uRW.GetByName": {
			Calls: func(i *mock.Interactor) {
				i.UserRW.EXPECT().GetByName(gomock.Any()).Return(nil, errors.New(""))
			}},
		"NO ERROR on return aRW.GetBySlug": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(nil, nil)
			}},
		"error on return BeforeCreationCheck": {
			Calls: func(i *mock.Interactor) {
				i.ArticleValidator.EXPECT().BeforeCreationCheck(gomock.Any()).Return(errors.New(""))
			}},
		"error on return Create": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().Create(gomock.Any()).Return(nil, errors.New(""))
			}},
		"error on tags add": {
			Calls: func(i *mock.Interactor) {
				i.TagsRW.EXPECT().Add(gomock.Any()).Return(errors.New("")).AnyTimes()
			}},
	}

	rick := testData.User("rick")
	slug := "article-title"

	// same as the happy case but with any parameter and called any number of times (including 0)
	validCalls := func(i *mock.Interactor) {
		i.UserRW.EXPECT().GetByName(gomock.Any()).Return(&rick, nil).AnyTimes()
		i.Slugger.EXPECT().NewSlug(gomock.Any()).Return(slug).AnyTimes()
		i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(nil, errors.New("not found")).AnyTimes()
		i.ArticleValidator.EXPECT().BeforeCreationCheck(gomock.Any()).Return(nil).AnyTimes()
		i.TagsRW.EXPECT().Add(gomock.Any()).AnyTimes()
		i.ArticleRW.EXPECT().Create(gomock.Any()).AnyTimes()
	}

	for testName, mutation := range mutations {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			i := mock.NewMockedInteractor(mockCtrl)
			mutation.Calls(&i) // put the tested call first (important)
			validCalls(&i)     // then fill the gaps with valid calls

			_, _, err := i.GetUCHandler().ArticlePost(rick.Name, domain.Article{})
			if mutation.ShouldPass {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}
