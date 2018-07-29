package uc_test

import (
	"testing"

	"github.com/err0r500/go-realworld-clean/domain"
	mock "github.com/err0r500/go-realworld-clean/implem/mock.uc"
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
		i := mock.NewMockedInteractor(mockCtrl)
		i.ArticleRW.EXPECT().GetRecentFiltered(filters).Return(expectedArticles, nil).Times(1)

		articles, count, err := i.GetUCHandler().GetArticles(10, offset, filters)
		assert.NoError(t, err)
		assert.Equal(t, 4, count)
		assert.Equal(t, expectedArticles[offset:], articles)
	})
}
