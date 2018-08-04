package uc

import (
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
	Create(username, email, password string) (*domain.User, error)
	GetByName(userName string) (*domain.User, error)
	GetByEmailAndPassword(email, password string) (*domain.User, error)
	Save(user domain.User) error
}

type ArticleRW interface {
	Create(domain.Article) (*domain.Article, error)
	Save(domain.Article) (*domain.Article, error)
	GetBySlug(slug string) (*domain.Article, error)
	GetByAuthorsNameOrderedByMostRecentAsc(usernames []string) ([]domain.Article, error)
	GetRecentFiltered(filters []domain.ArticleFilter) ([]domain.Article, error)
	Delete(slug string) error
}

type CommentRW interface {
	Create(comment domain.Comment) (*domain.Comment, error)
	GetByID(id int) (*domain.Comment, error)
	Delete(id int) error
}

type TagsRW interface {
	GetAll() ([]string, error)
	Add(newTags []string) error
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
