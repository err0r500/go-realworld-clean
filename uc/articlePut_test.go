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

func TestInteractor_ArticlePut_happycase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	jane := testData.User("jane")
	origArticle := testData.Article("janeArticle")
	req := domain.Article{
		Title:       "newTitle",
		Description: "newDescription",
	}

	toInsert := origArticle
	toInsert.Title = req.Title
	toInsert.Description = req.Description

	i := mock.NewMockedInteractor(mockCtrl)
	i.ArticleRW.EXPECT().GetBySlug(origArticle.Slug).Return(&origArticle, nil).Times(1)
	i.ArticleValidator.EXPECT().BeforeUpdateCheck(&toInsert).Return(nil).Times(1)
	i.ArticleRW.EXPECT().Save(toInsert)
	i.UserRW.EXPECT().GetByName(jane.Name).Return(&jane, nil).Times(1)

	_, _, err := i.GetUCHandler().ArticlePut(jane.Name, origArticle.Slug, map[domain.ArticleUpdatableField]*string{domain.Title: &req.Title, domain.Description: &req.Description})
	assert.NoError(t, err)
}

func TestInteractor_ArticlePut_fails(t *testing.T) {
	mutations := map[string]mock.Tester{
		"shouldPass": {
			Calls: func(i *mock.Interactor) { // change nothing
			},
			ShouldPass: true},
		"error return on getBySlug": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(nil, errors.New(""))
			}},
		"nil return on getBySlug": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(nil, nil)
			}},
		"error return on articleValidator": {
			Calls: func(i *mock.Interactor) {
				i.ArticleValidator.EXPECT().BeforeUpdateCheck(gomock.Any()).Return(errors.New(""))
			}},
		"error return on uRW.GetByName": {
			Calls: func(i *mock.Interactor) {
				i.UserRW.EXPECT().GetByName(gomock.Any()).Return(nil, errors.New(""))
			}},
		"error return on aRW.Save": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().Save(gomock.Any()).Return(nil, errors.New(""))
			}},
	}

	jane := testData.User("jane")
	origArticle := testData.Article("janeArticle")

	// same as the happy case but with any parameter and called any number of times (including 0)
	validCalls := func(i *mock.Interactor) {
		i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(&origArticle, nil).AnyTimes()
		i.ArticleValidator.EXPECT().BeforeUpdateCheck(gomock.Any()).Return(nil).AnyTimes()
		i.ArticleRW.EXPECT().Save(gomock.Any()).AnyTimes()
		i.UserRW.EXPECT().GetByName(gomock.Any()).Return(&jane, nil).AnyTimes()
	}

	for testName, mutation := range mutations {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			i := mock.NewMockedInteractor(mockCtrl)
			mutation.Calls(&i) // put the tested call first (important)
			validCalls(&i)     // then fill the gaps with valid calls

			_, _, err := i.GetUCHandler().ArticlePut(jane.Name, origArticle.Slug, map[domain.ArticleUpdatableField]*string{})
			if mutation.ShouldPass {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}
