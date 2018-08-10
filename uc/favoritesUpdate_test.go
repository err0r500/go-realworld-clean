package uc_test

import (
	"errors"
	"testing"

	"github.com/err0r500/go-realworld-clean/domain"
	mock "github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInteractor_FavoritesUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := testData.User("rick")
	article := testData.Article("jane")
	articleWithoutFav := article
	articleWithoutFav.FavoritedBy = []domain.User{}

	i := mock.NewMockedInteractor(mockCtrl)
	i.ArticleRW.EXPECT().GetBySlug(article.Slug).Return(&articleWithoutFav, nil)
	i.UserRW.EXPECT().GetByName(user.Name).Return(&user, nil)
	i.ArticleRW.EXPECT().Save(article).Return(&article, nil)

	returnedUser, returnedArticle, err := i.GetUCHandler().FavoritesUpdate(user.Name, article.Slug, true)
	assert.NoError(t, err)
	assert.Equal(t, user, *returnedUser)
	assert.Equal(t, article, *returnedArticle)
}

func TestInteractor_FavoritesUpdate_fails(t *testing.T) {
	mutations := map[string]mock.Tester{
		"shouldPass": {
			Calls: func(i *mock.Interactor) { // change nothing
			},
			ShouldPass: true},
		"failed get user by name": {
			Calls: func(i *mock.Interactor) {
				i.UserRW.EXPECT().GetByName(gomock.Any()).Return(nil, errors.New(""))
			},
		},
		"failed get by slug": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(nil, errors.New(""))
			},
		},
		"failed save article": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().Save(gomock.Any()).Return(nil, errors.New(""))
			},
		},
	}

	user := testData.User("rick")
	article := testData.Article("jane")
	articleWithoutFav := article
	articleWithoutFav.FavoritedBy = []domain.User{}

	// same as the happy case but with any parameter and called any number of times (including 0)
	validCalls := func(i *mock.Interactor) {
		i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(&articleWithoutFav, nil).AnyTimes()
		i.UserRW.EXPECT().GetByName(gomock.Any()).Return(&user, nil).AnyTimes()
		i.ArticleRW.EXPECT().Save(gomock.Any()).Return(&article, nil).AnyTimes()
	}

	for testName, mutation := range mutations {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			i := mock.NewMockedInteractor(mockCtrl)
			mutation.Calls(&i) // put the tested call first (important)
			validCalls(&i)     // then fill the gaps with valid calls

			_, _, err := i.GetUCHandler().FavoritesUpdate(user.Name, article.Slug, true)
			if mutation.ShouldPass {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}
