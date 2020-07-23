package uc

import (
	"context"
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
	ProfileGet(requestingUserName, userName string) (profile *domain.User, follows bool, err error)
	ProfileUpdateFollow(loggedInUsername, username string, follow bool) (user *domain.User, err error)
}

type UserLogic interface {
	UserCreate(username, email, password string) (user *domain.User, token string, err error)
	UserLogin(email, password string) (user *domain.User, token string, err error)
	UserGet(userName string) (user *domain.User, token string, err error)
	UserEdit(userName string, fieldsToUpdate map[domain.UserUpdatableProperty]*string) (user *domain.User, token string, err error)
}

type ArticlesLogic interface {
	ArticlesFeed(ctx context.Context, username string, limit, offset int) (requestingUser *domain.User, articles domain.ArticleCollection, totalArticleCount int, err error)
	GetArticles(ctx context.Context, username string, limit, offset int, filters []domain.ArticleFilter) (requestingUser *domain.User, articles domain.ArticleCollection, totalArticleCount int, err error)
}

type ArticleLogic interface {
	ArticleGet(ctx context.Context, slug, username string) (*domain.User, *domain.Article, error)
	ArticlePost(ctx context.Context, username string, article domain.Article) (*domain.User, *domain.Article, error)
	ArticlePut(ctx context.Context, username, slug string, fieldsToUpdate map[domain.ArticleUpdatableField]*string) (*domain.User, *domain.Article, error)
	ArticleDelete(ctx context.Context, username, slug string) error
}

type CommentsLogic interface {
	CommentsGet(ctx context.Context, slug string) ([]domain.Comment, error)
	CommentsPost(ctx context.Context, username, slug, comment string) (*domain.Comment, error)
	CommentsDelete(ctx context.Context, username, slug string, id int) error
}

type FavoritesLogic interface {
	FavoritesUpdate(ctx context.Context, username, slug string, favortie bool) (*domain.User, *domain.Article, error)
}

type TagsLogic interface {
	Tags() ([]string, error)
}

type HandlerConstructor struct {
	Logger           Logger
	UserRW           UserRW
	ArticleRW        ArticleRW
	CommentRW        CommentRW
	UserValidator    UserValidator
	AuthHandler      AuthHandler
	Slugger          Slugger
	ArticleValidator ArticleValidator
	TagsRW           TagsRW
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
	if c.TagsRW == nil {
		log.Fatal("missing TagsRW")
	}
	if c.CommentRW == nil {
		log.Fatal("missing CommentRW")
	}

	return interactor{
		logger:           c.Logger,
		userRW:           c.UserRW,
		articleRW:        c.ArticleRW,
		userValidator:    c.UserValidator,
		authHandler:      c.AuthHandler,
		slugger:          c.Slugger,
		articleValidator: c.ArticleValidator,
		tagsRW:           c.TagsRW,
		commentRW:        c.CommentRW,
	}
}
