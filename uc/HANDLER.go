package uc

import (
	"log"

	"github.com/err0r500/go-realworld-clean/domain"
)

type Handler interface {
	ProfileLogic
	UserLogic
	ArticlesLogic
	ArticleLogic
	CommentsLogic
	FavoritesLogic
	TagsLogic
}

type ProfileLogic interface {
	ProfileGet(userName string) (profile *domain.Profile, err error)
	ProfileUpdateFollow(loggedInUsername, username string, follow bool) (user *domain.User, err error)
}

type UserLogic interface {
	UserCreate(username, email, password string) (user *domain.User, token string, err error)
	UserLogin(email, password string) (user *domain.User, token string, err error)
	UserGet(userName string) (user *domain.User, token string, err error)
	UserEdit(userName string, newUser map[UpdatableProperty]*string) (user *domain.User, token string, err error)
}

type ArticlesLogic interface {
	ArticlesFeed(username string, limit, offset int) (articles domain.ArticleCollection, totalArticleCount int, err error)
	GetArticles(limit, offset int, filters Filters) (articles domain.ArticleCollection, totalArticleCount int, err error)
}

type ArticleLogic interface {
	ArticleGet(slug string) (*domain.Article, error)
	ArticlePost(username string, article domain.Article) (*domain.Article, error)
	ArticlePut(username, slug string, article domain.Article) (*domain.Article, error)
	ArticleDelete(username, slug string) error
}

type CommentsLogic interface {
	CommentsGet(slug string) ([]domain.Comment, error)
	CommentsPost(username, slug, comment string) (*domain.Comment, error)
	CommentsDelete(username, slug, id string) error
}

type FavoritesLogic interface {
	FavoritesUpdate(username, slug string, favortie bool) (*domain.Article, error)
}

type TagsLogic interface {
	Tags() ([]string, error)
}

type HandlerConstructor struct {
	Logger           Logger
	UserRW           UserRW
	ArticleRW        ArticleRW
	UserValidator    UserValidator
	AuthHandler      AuthHandler
	Slugger          Slugger
	ArticleValidator ArticleValidator
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
	if c.Slugger == nil {
		log.Fatal("missing Slugger")
	}
	if c.ArticleValidator == nil {
		log.Fatal("missing ArticleValidator")
	}

	return interactor{
		logger:           c.Logger,
		userRW:           c.UserRW,
		articleRW:        c.ArticleRW,
		userValidator:    c.UserValidator,
		authHandler:      c.AuthHandler,
		slugger:          c.Slugger,
		articleValidator: c.ArticleValidator,
	}
}
