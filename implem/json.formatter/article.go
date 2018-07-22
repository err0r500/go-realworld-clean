package formatter

import "github.com/err0r500/go-realworld-clean/domain"

type Article struct {
	Title          string   `json:"title"`
	Slug           string   `json:"slug"`
	Description    string   `json:"description"`
	Body           string   `json:"body"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	Author         Profile  `json:"author"`
	Tags           []string `json:"tagList"`
	Favorite       bool     `json:"favorited"`
	FavoritesCount int      `json:"favoritesCount"`
}

func NewArticleFromDomain(article domain.Article, isFollowingAuthor bool) Article {
	return Article{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		CreatedAt:      article.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt:      article.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		Author:         NewProfileFromDomain(article.Author, isFollowingAuthor),
		Tags:           article.TagList,
		Favorite:       article.Favorited,
		FavoritesCount: article.FavoritesCount,
	}
}
