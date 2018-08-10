// +build !netgo

package mock

import (
	"log"

	"github.com/err0r500/go-realworld-clean/uc"
	"github.com/golang/mock/gomock"
)

// Interactor : is used in order to update its properties accordingly to each test conditions
type Interactor struct {
	Logger           *MockLogger
	UserRW           *MockUserRW
	ArticleRW        *MockArticleRW
	UserValidator    *MockUserValidator
	AuthHandler      *MockAuthHandler
	Slugger          *MockSlugger
	ArticleValidator *MockArticleValidator
	TagsRW           *MockTagsRW
	CommentRW        *MockCommentRW
}

type Tester struct {
	Calls      func(*Interactor)
	ShouldPass bool
}

type SimpleLogger struct{}

func (SimpleLogger) Log(logs ...interface{}) {
	log.Println(logs...)
}

//NewMockedInteractor : the Interactor constructor
func NewMockedInteractor(mockCtrl *gomock.Controller) Interactor {
	return Interactor{
		Logger:           NewMockLogger(mockCtrl),
		UserRW:           NewMockUserRW(mockCtrl),
		ArticleRW:        NewMockArticleRW(mockCtrl),
		UserValidator:    NewMockUserValidator(mockCtrl),
		AuthHandler:      NewMockAuthHandler(mockCtrl),
		Slugger:          NewMockSlugger(mockCtrl),
		ArticleValidator: NewMockArticleValidator(mockCtrl),
		TagsRW:           NewMockTagsRW(mockCtrl),
		CommentRW:        NewMockCommentRW(mockCtrl),
	}
}

//GetUCHandler : returns a uc.interactor in order to call its methods aka the use cases to test
func (i Interactor) GetUCHandler() uc.Handler {
	return uc.HandlerConstructor{
		Logger:           i.Logger,
		UserRW:           i.UserRW,
		ArticleRW:        i.ArticleRW,
		UserValidator:    i.UserValidator,
		AuthHandler:      i.AuthHandler,
		Slugger:          i.Slugger,
		ArticleValidator: i.ArticleValidator,
		TagsRW:           i.TagsRW,
		CommentRW:        i.CommentRW,
	}.New()
}
