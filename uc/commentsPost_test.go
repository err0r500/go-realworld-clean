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

func TestInteractor_CommentsPost_happycase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := testData.User("rick")
	article := testData.Article("jane")
	articleWithoutComment := article
	articleWithoutComment.Comments = []domain.Comment{} //reset comments
	comment := article.Comments[0]

	i := mock.NewMockedInteractor(mockCtrl)
	i.UserRW.EXPECT().GetByName(user.Name).Return(&user, nil)
	i.ArticleRW.EXPECT().GetBySlug(article.Slug).Return(&articleWithoutComment, nil)
	i.CommentRW.EXPECT().Create(domain.Comment{Body: comment.Body, Author: user}).Return(&comment, nil)
	i.ArticleRW.EXPECT().Save(article)

	returnedComment, err := i.GetUCHandler().CommentsPost(user.Name, article.Slug, comment.Body)
	assert.NoError(t, err)
	assert.Equal(t, comment, *returnedComment)
}

func TestInteractor_CommentsPost_fails(t *testing.T) {
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
		"failed get create comment": {
			Calls: func(i *mock.Interactor) {
				i.CommentRW.EXPECT().Create(gomock.Any()).Return(nil, errors.New(""))
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
	articleWithoutComment := article
	articleWithoutComment.Comments = []domain.Comment{} //reset comments
	comment := article.Comments[0]

	// same as the happy case but with any parameter and called any number of times (including 0)
	validCalls := func(i *mock.Interactor) {
		i.UserRW.EXPECT().GetByName(gomock.Any()).Return(&user, nil).AnyTimes()
		i.ArticleRW.EXPECT().GetBySlug(gomock.Any()).Return(&articleWithoutComment, nil).AnyTimes()
		i.CommentRW.EXPECT().Create(gomock.Any()).Return(&comment, nil).AnyTimes()
		i.ArticleRW.EXPECT().Save(gomock.Any()).AnyTimes()
	}

	for testName, mutation := range mutations {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			i := mock.NewMockedInteractor(mockCtrl)
			mutation.Calls(&i) // put the tested call first (important)
			validCalls(&i)     // then fill the gaps with valid calls

			_, err := i.GetUCHandler().CommentsPost(user.Name, article.Slug, comment.Body)
			if mutation.ShouldPass {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}
