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

var art1 = domain.Article{Slug: "1"}
var art2 = domain.Article{Slug: "2"}
var art3 = domain.Article{Slug: "3"}
var art4 = domain.Article{Slug: "4"}

func TestInteractor_ArticlesFeed_happycases(t *testing.T) {
	rick := testData.User("rick")
	rick.FollowIDs = []string{testData.User("jane").Name}
	expectedArticles := domain.ArticleCollection{art1, art2, art3, art4}

	t.Run("most obvious", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		i := mock.NewMockedInteractor(mockCtrl)
		i.UserRW.EXPECT().GetByName(rick.Name).Return(&rick, nil)
		i.ArticleRW.EXPECT().GetByAuthorsNameOrderedByMostRecentAsc(rick.FollowIDs).Return(expectedArticles, nil).Times(1)

		user, articles, count, err := i.GetUCHandler().ArticlesFeed(rick.Name, 4, 0)
		assert.NoError(t, err)
		assert.Equal(t, 4, count)
		assert.Equal(t, expectedArticles, articles)
		assert.Equal(t, rick, *user)
	})

	t.Run("total count", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		i := mock.NewMockedInteractor(mockCtrl)
		i.UserRW.EXPECT().GetByName(rick.Name).Return(&rick, nil)
		i.ArticleRW.EXPECT().GetByAuthorsNameOrderedByMostRecentAsc(rick.FollowIDs).Return(expectedArticles, nil).Times(1)

		user, articles, count, err := i.GetUCHandler().ArticlesFeed(rick.Name, 2, 0)
		assert.NoError(t, err)
		assert.Equal(t, 4, count)
		assert.Equal(t, expectedArticles[:2], articles)
		assert.Equal(t, rick, *user)
	})

	t.Run("empty count", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		i := mock.NewMockedInteractor(mockCtrl)
		user, _, count, err := i.GetUCHandler().ArticlesFeed(rick.Name, -2, 0)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
		assert.Nil(t, user)
	})
}

func TestInteractor_ArticlesFeed_fails(t *testing.T) {
	mutations := map[string]mock.Tester{
		"shouldPass": {
			Calls: func(i *mock.Interactor) { // change nothing
			},
			ShouldPass: true},
		"error return on uRW.GetByName": {
			Calls: func(i *mock.Interactor) {
				i.UserRW.EXPECT().GetByName(gomock.Any()).Return(nil, errors.New(""))
			}},
		"error return on GetByAuthorsNameOrderedByMostRecentAsc": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().GetByAuthorsNameOrderedByMostRecentAsc(gomock.Any()).Return(nil, errors.New(""))
			}},
	}
	rick := testData.User("rick")
	rick.FollowIDs = []string{testData.User("jane").Name}
	expectedArticles := domain.ArticleCollection{art1, art2, art3, art4}

	// same as the happy case but with any parameter and called any number of times (including 0)
	validCalls := func(i *mock.Interactor) {
		i.UserRW.EXPECT().GetByName(gomock.Any()).Return(&rick, nil).AnyTimes()
		i.ArticleRW.EXPECT().GetByAuthorsNameOrderedByMostRecentAsc(gomock.Any()).Return(expectedArticles, nil).AnyTimes()
	}

	for testName, mutation := range mutations {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			i := mock.NewMockedInteractor(mockCtrl)
			mutation.Calls(&i) // put the tested call first (important)
			validCalls(&i)     // then fill the gaps with valid calls

			_, _, _, err := i.GetUCHandler().ArticlesFeed(rick.Name, 2, 0)
			if mutation.ShouldPass {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}
