package domain

import (
	"time"
)

type Article struct {
	Slug        string
	Title       string
	Description string
	Body        string
	TagList     []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FavoritedBy []User
	Author      User
	Comments    []Comment
}

type Comment struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
	Author    User
}

type ArticleUpdatableField int

const (
	Title ArticleUpdatableField = iota
	Description
	Body
)

func UpdateArticle(initial *Article, opts ...func(fields *Article)) {
	for _, v := range opts {
		v(initial)
	}
}

func SetArticleTitle(input *string) func(fields *Article) {
	return func(initial *Article) {
		if input != nil {
			initial.Title = *input
		}
	}
}

func SetArticleDescription(input *string) func(fields *Article) {
	return func(initial *Article) {
		if input != nil {
			initial.Description = *input
		}
	}
}

func SetArticleBody(input *string) func(fields *Article) {
	return func(initial *Article) {
		if input != nil {
			initial.Body = *input
		}
	}
}

type ArticleFilter func(Article) bool

func ArticleHasTag(tag string) ArticleFilter {
	return func(article Article) bool {
		for _, articleTag := range article.TagList {
			if articleTag == tag {
				return true
			}
		}
		return false
	}
}

func ArticleHasAuthor(authorName string) ArticleFilter {
	return func(article Article) bool {
		return article.Author.Name == authorName
	}
}

func ArticleIsFavoritedBy(username string) ArticleFilter {
	return func(article Article) bool {
		if username == "" {
			return false
		}
		for _, user := range article.FavoritedBy {
			if user.Name == username {
				return true
			}
		}
		return false
	}
}

type ArticleCollection []Article

func (articles ArticleCollection) ApplyLimitAndOffset(limit, offset int) ArticleCollection {
	if limit <= 0 {
		return []Article{}
	}

	articlesSize := len(articles)
	min := offset
	if min < 0 {
		min = 0
	}

	if min > articlesSize {
		return []Article{}
	}

	max := min + limit
	if max > articlesSize {
		max = articlesSize
	}

	return articles[min:max]
}

func (article *Article) UpdateComments(comment Comment, add bool) {
	if add {
		article.Comments = append(article.Comments, comment)
		return
	}

	for i := 0; i < len(article.Comments); i++ {
		if article.Comments[i].ID == comment.ID {
			article.Comments = append(article.Comments[:i], article.Comments[i+1:]...) // memory leak ? https://github.com/golang/go/wiki/SliceTricks
		}
	}
}

func (article *Article) UpdateFavoritedBy(user User, add bool) {
	if add {
		article.FavoritedBy = append(article.FavoritedBy, user)
		return
	}

	for i := 0; i < len(article.FavoritedBy); i++ {
		if article.FavoritedBy[i].Name == user.Name {
			article.FavoritedBy = append(article.FavoritedBy[:i], article.FavoritedBy[i+1:]...) // memory leak ? https://github.com/golang/go/wiki/SliceTricks
		}
	}
}
