package formatter

import (
	"github.com/err0r500/go-realworld-clean/domain"
)

const dateLayout = "2006-01-02T15:04:05.999Z"

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

func NewArticleFromDomain(article domain.Article, user *domain.User) Article {
	isFollowingAuthor := false
	favorite := false
	if user != nil {
		for _, userName := range user.FollowIDs {
			if userName == article.Author.Name {
				isFollowingAuthor = true
				break
			}
		}

		favorite = domain.ArticleIsFavoritedBy(user.Name)(article)
	}

	return Article{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		CreatedAt:      article.CreatedAt.UTC().Format(dateLayout),
		UpdatedAt:      article.UpdatedAt.UTC().Format(dateLayout),
		Author:         NewProfileFromDomain(article.Author, isFollowingAuthor),
		Tags:           article.TagList,
		Favorite:       favorite,
		FavoritesCount: len(article.FavoritedBy),
	}
}

func NewArticlesFromDomain(user *domain.User, articles ...domain.Article) []Article {
	ret := []Article{} // return at least an empty array (not nil)

	for _, article := range articles {
		ret = append(ret, NewArticleFromDomain(article, user))
	}

	return ret
}
