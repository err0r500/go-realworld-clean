package uc

import (
	"context"

	"github.com/err0r500/go-realworld-clean/domain"
)

// interactor : the struct that will have as properties all the IMPLEMENTED interfaces
// in order to provide them to its methods : the use cases and implement the Handler interface
type interactor struct {
	logger           Logger
	userRW           UserRW
	articleRW        ArticleRW
	userValidator    UserValidator
	articleValidator ArticleValidator
	authHandler      AuthHandler
	slugger          Slugger
	commentRW        CommentRW
	tagsRW           TagsRW
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
	Create(ctx context.Context, username, email, password string) (*domain.User, bool)
	GetByName(ctx context.Context, userName string) (*domain.User, bool)
	GetByEmailAndPassword(ctx context.Context, email, password string) (*domain.User, bool)
	Save(ctx context.Context, user domain.User) bool
}

type ArticleRW interface {
	Create(ctx context.Context, _ domain.Article) (*domain.Article, bool)
	Save(ctx context.Context, _ domain.Article) (*domain.Article, bool)
	GetBySlug(ctx context.Context, slug string) (*domain.Article, bool)
	GetByAuthorsNameOrderedByMostRecentAsc(ctx context.Context, usernames []string) ([]domain.Article, bool)
	GetRecentFiltered(ctx context.Context, filters []domain.ArticleFilter) ([]domain.Article, bool)
	Delete(ctx context.Context, slug string) bool
}

type CommentRW interface {
	Create(ctx context.Context, comment domain.Comment) (*domain.Comment, bool)
	GetByID(ctx context.Context, id int) (*domain.Comment, bool)
	Delete(ctx context.Context, id int) bool
}

type TagsRW interface {
	GetAll(ctx context.Context) ([]string, bool)
	Add(ctx context.Context, newTags []string) bool
}

type Slugger interface {
	NewSlug(string) string
}

type UserValidator interface {
	CheckUser(user domain.User) error
}

type ArticleValidator interface {
	BeforeCreationCheck(article *domain.Article) error
	BeforeUpdateCheck(article *domain.Article) error
}
