package uc

import (
	"log"

	"github.com/err0r500/go-realworld-clean/domain"
)

type Handler interface {
	ProfileGet(userName string) (profile *domain.Profile, err error)
	ProfileUpdateFollow(loggedInuserName, username string, follow bool) (user *domain.User, err error)

	UserCreate(username, email, password string) (user *domain.User, token string, err error)
	UserLogin(email, password string) (user *domain.User, token string, err error)
	UserGet(userName string) (user *domain.User, token string, err error)
	UserEdit(userName string, newUser map[UpdatableProperty]*string) (user *domain.User, token string, err error)

	ArticlesFeed(username string, limit, offset int) (articles domain.ArticleCollection, totalArticleCount int, err error)
	GetArticles(limit, offset int, filters Filters) (articles domain.ArticleCollection, totalArticleCount int, err error)

	ArticleGet(slug string) (*domain.Article, error)
	ArticlePost(username string, article domain.Article) (*domain.Article, error)
	ArticlePut(username string, slug string, article domain.Article) (*domain.Article, error)
	ArticleDelete(username string, slug string) error
}

type HandlerConstructor struct {
	Logger        Logger
	UserRW        UserRW
	ArticleRW     ArticleRW
	UserValidator UserValidator
	AuthHandler   AuthHandler
}

func (c HandlerConstructor) New() Handler {
	if c.Logger == nil {
		log.Fatal("missing Logger")
	}
	if c.UserRW == nil {
		log.Fatal("missing UserRW")
	}
	if c.ArticleRW == nil {
		log.Fatal("missing ArticleRW")
	}
	if c.UserValidator == nil {
		log.Fatal("missing UserValidator")
	}
	if c.AuthHandler == nil {
		log.Fatal("missing AuthHandler")
	}

	return interactor{
		logger:        c.Logger,
		userRW:        c.UserRW,
		articleRW:     c.ArticleRW,
		userValidator: c.UserValidator,
		authHandler:   c.AuthHandler,
	}
}

// NewHandler : the interactor constructor, use this in order to avoid null pointers at runtime
//func NewHandler(Logger Logger, uRW UserRW, arw ArticleRW, validator UserValidator, handler AuthHandler) Handler {
//	return interactor{
//		Logger:        Logger,
//		UserRW:        uRW,
//		ArticleRW:     arw,
//		UserValidator: validator,
//		AuthHandler:   handler,
//	}
//}

// interactor : the struct that will have as properties all the IMPLEMENTED interfaces
// in order to provide them to its methods : the use cases
type interactor struct {
	logger        Logger
	userRW        UserRW
	articleRW     ArticleRW
	userValidator UserValidator
	authHandler   AuthHandler
}

// Logger : only used to log stuff
type Logger interface {
	Log(...interface{})
}

type AuthHandler interface {
	GenUserToken(userName string) (token string, err error)
	GetUserName(token string) (userName string, err error)
}

type UserRW interface {
	Create(username, email, password string) (*domain.User, error)
	GetByName(userName string) (*domain.User, error)
	GetByEmailAndPassword(email, password string) (*domain.User, error)
	Save(user domain.User) error
}

type ArticleRW interface {
	GetByAuthorsNameOrderedByMostRecentAsc(usernames []string) ([]domain.Article, error)
	GetRecentFiltered(filters Filters) ([]domain.Article, error)
}

type UserValidator interface {
	CheckUser(user domain.User) error
}
