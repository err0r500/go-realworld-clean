// +build !netgo

package uc

import (
	"log"

	"github.com/err0r500/go-realworld-clean/uc"
	"github.com/golang/mock/gomock"
)

// MockedInteractor : is used in order to update its properties accordingly to each test conditions
type MockedInteractor struct {
	Logger        uc.Logger
	UserRW        *MockUserRW
	ArticleRW     *MockArticleRW
	UserValidator *MockUserValidator
	AuthHandler   *MockAuthHandler
}

type SimpleLogger struct{}

func (SimpleLogger) Log(logs ...interface{}) {
	log.Println(logs)
}

//NewMockedInteractor : the MockedInteractor constructor
func NewMockedInteractor(mockCtrl *gomock.Controller) MockedInteractor {
	return MockedInteractor{
		Logger:        SimpleLogger{},
		UserRW:        NewMockUserRW(mockCtrl),
		ArticleRW:     NewMockArticleRW(mockCtrl),
		UserValidator: NewMockUserValidator(mockCtrl),
		AuthHandler:   NewMockAuthHandler(mockCtrl),
	}
}

//GetUCHandler : returns a uc.interactor in order to call its methods aka the use cases to test
func (i MockedInteractor) GetUCHandler() uc.Handler {
	return uc.HandlerConstructor{
		Logger:        i.Logger,
		UserRW:        i.UserRW,
		ArticleRW:     i.ArticleRW,
		UserValidator: i.UserValidator,
		AuthHandler:   i.AuthHandler,
	}.New()
}
