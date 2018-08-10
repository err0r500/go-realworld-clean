package uc_test

import (
	"testing"

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
		filters := uc.NewFilters("jane", "", "")
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
}
