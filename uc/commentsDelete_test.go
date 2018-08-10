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

func TestInteractor_CommentsDelete_happycase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := testData.User("rick")
	article := testData.Article("jane")
	comment := article.Comments[0]
	articleWithoutTheComment := article
	articleWithoutTheComment.Comments = []domain.Comment{}

	i := mock.NewMockedInteractor(mockCtrl)
	i.CommentRW.EXPECT().GetByID(comment.ID).Return(&comment, nil)
	i.CommentRW.EXPECT().Delete(comment.ID).Return(nil)
	i.ArticleRW.EXPECT().GetBySlug(article.Slug).Return(&article, nil)
	i.ArticleRW.EXPECT().Save(articleWithoutTheComment)

	assert.NoError(t, i.GetUCHandler().CommentsDelete(user.Name, article.Slug, comment.ID))
}

func TestInteractor_CommentsDelete_fails(t *testing.T) {
	mutations := map[string]mock.Tester{
		"shouldPass": {
			Calls: func(i *mock.Interactor) { // change nothing
			},
			ShouldPass: true},
		"failed to get comment by ID": {
			Calls: func(i *mock.Interactor) {
				i.CommentRW.EXPECT().GetByID(gomock.Any()).Return(nil, errors.New(""))
			}},
		"returned a comment belonging to another user": {
			Calls: func(i *mock.Interactor) {
				i.CommentRW.EXPECT().GetByID(gomock.Any()).Return(&domain.Comment{Author: domain.User{Name: "hey"}}, nil)
			}},
		"failed to delete comment": {
			Calls: func(i *mock.Interactor) {
				i.CommentRW.EXPECT().Delete(gomock.Any()).Return(errors.New(""))
			}},
		"failed to get article by slug comment": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(nil, errors.New(""))
			}},
		"failed to save article": {
			Calls: func(i *mock.Interactor) {
				i.ArticleRW.EXPECT().Save(gomock.Any()).Return(nil, errors.New(""))
			}},
	}

	user := testData.User("rick")
	article := testData.Article("jane")
	comment := article.Comments[0]
	articleWithoutTheComment := article
	articleWithoutTheComment.Comments = []domain.Comment{}

	// same as the happy case but with any parameter and called any number of times (including 0)
	validCalls := func(i *mock.Interactor) {
		i.CommentRW.EXPECT().GetByID(gomock.Any()).Return(&comment, nil).AnyTimes()
		i.CommentRW.EXPECT().Delete(gomock.Any()).Return(nil).AnyTimes()
		i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(&article, nil).AnyTimes()
		i.ArticleRW.EXPECT().Save(gomock.Any()).AnyTimes()

	}

	for testName, mutation := range mutations {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			i := mock.NewMockedInteractor(mockCtrl)
			mutation.Calls(&i) // put the tested call first (important)
			validCalls(&i)     // then fill the gaps with valid calls

			err := i.GetUCHandler().CommentsDelete(user.Name, article.Slug, comment.ID)
			if mutation.ShouldPass {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}
