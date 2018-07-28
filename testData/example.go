// +build !netgo

package testData

import (
	"time"

	"github.com/err0r500/go-realworld-clean/domain"
)

var rickBio = "Rick biography string"
var janeImg = "jane img link"

func User(name string) domain.User {
	switch name {
	case "rick":
		return rick
	default:
		return jane
	}
}
func Article(name string) domain.Article {
	switch name {
	default:
		return janeArticle
	}
}

var rick = domain.User{
	Name:      "rick",
	Email:     "rick@example.com",
	Bio:       &rickBio,
	ImageLink: nil,
	Password:  "rickPassword",
}

var jane = domain.User{
	Name:      "jane",
	Email:     "jane@example.com",
	Bio:       nil,
	ImageLink: &janeImg,
	Password:  "janePassword",
}

func Profile(name string) domain.Profile {
	switch name {
	case "janeFollowingRick":
		return janeFollowingRick
	default:
		return janeFollowingRick
	}
}

var janeFollowingRick = domain.Profile{
	User:      rick,
	Following: true,
}

var janeArticle = domain.Article{
	Slug:           "articleSlug",
	Title:          "articleTitle",
	Description:    "description",
	Body:           "body",
	TagList:        []string{"tagList"},
	CreatedAt:      time.Now(),
	UpdatedAt:      time.Now(),
	Favorited:      true,
	FavoritesCount: 123,
	Author:         domain.Profile{User: jane, Following: false},
	Comments: []domain.Comment{
		{ID: "1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Body:      "commentBody",
			Author:    domain.Profile{User: rick, Following: false},
		},
	},
}
