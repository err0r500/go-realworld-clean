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
	Author         User      `json:"author"`
	Comments       []Comment `json:"comments"`
}

type Comment struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Author    User      `json:"author"`
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
