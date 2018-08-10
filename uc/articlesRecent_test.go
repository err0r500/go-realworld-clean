package uc_test

import (
	"testing"

	"errors"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/err0r500/go-realworld-clean/uc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInteractor_GetArticles(t *testing.T) {
	expectedArticles := domain.ArticleCollection{art1, art2, art3, art4}

	t.Run("most obvious", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		offset := 2
		filters := uc.NewFilters("jane", "bla", "true")
		testUser := testData.User("rick")
		i := mock.NewMockedInteractor(mockCtrl)
		i.ArticleRW.EXPECT().GetRecentFiltered(filters).Return(expectedArticles, nil).Times(1)
		i.UserRW.EXPECT().GetByName(testUser.Name).Return(&testUser, nil).Times(1)

		user, articles, count, err := i.GetUCHandler().GetArticles(testUser.Name, 10, offset, filters)
		assert.NoError(t, err)
		assert.Equal(t, 4, count)
		assert.Equal(t, expectedArticles[offset:], articles)
		assert.Equal(t, testData.User("rick"), *user)
	})

	t.Run("empty", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		testUser := testData.User("rick")
		i := mock.NewMockedInteractor(mockCtrl)

		user, _, count, err := i.GetUCHandler().GetArticles(testUser.Name, -10, 2, nil)
		assert.NoError(t, err)
		assert.Nil(t, user)
		assert.Equal(t, 0, count)
	})
}

func TestInteractor_ArticlesRecent_fails(t *testing.T) {
	expectedArticles := domain.ArticleCollection{art1, art2, art3, art4}

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
				i.ArticleRW.EXPECT().GetRecentFiltered(gomock.Any()).Return(nil, errors.New(""))
			}},
	}
	testUser := testData.User("rick")

	// same as the happy case but with any parameter and called any number of times (including 0)
	validCalls := func(i *mock.Interactor) {
		i.ArticleRW.EXPECT().GetRecentFiltered(gomock.Any()).Return(expectedArticles, nil).AnyTimes()
		i.UserRW.EXPECT().GetByName(gomock.Any()).Return(&testUser, nil).AnyTimes()
	}

	for testName, mutation := range mutations {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			i := mock.NewMockedInteractor(mockCtrl)
			mutation.Calls(&i) // put the tested call first (important)
			validCalls(&i)     // then fill the gaps with valid calls

			_, _, _, err := i.GetUCHandler().GetArticles(testUser.Name, 10, 2, nil)
			if mutation.ShouldPass {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}
