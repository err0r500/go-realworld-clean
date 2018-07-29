package domain

import (
	"time"
)

type Article struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         Profile   `json:"author"`
	Comments       []Comment `json:"comments"`
}

type Comment struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Author    Profile   `json:"author"`
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
