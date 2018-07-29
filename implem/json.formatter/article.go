package formatter

import "github.com/err0r500/go-realworld-clean/domain"

const dateFormat = "2006-01-02T15:04:05.999Z"

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

func NewArticleFromDomain(article domain.Article) Article {
	return Article{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		CreatedAt:      article.CreatedAt.UTC().Format(dateFormat),
		UpdatedAt:      article.UpdatedAt.UTC().Format(dateFormat),
		Author:         NewProfileFromDomain(article.Author, false), //fixme : check this !
		Tags:           article.TagList,
		Favorite:       article.Favorited,
		FavoritesCount: article.FavoritesCount,
	}
}

func NewArticlesFromDomain(articles ...domain.Article) []Article {
	var ret []Article
	for _, article := range articles {
		ret = append(ret, NewArticleFromDomain(article))
	}

	return ret
}
