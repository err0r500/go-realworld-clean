package domain_test

import (
	"testing"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/stretchr/testify/assert"
)

var art1 = domain.Article{Slug: "1"}
var art2 = domain.Article{Slug: "2"}
var art3 = domain.Article{Slug: "3"}
var art4 = domain.Article{Slug: "4"}

func TestArticleCollection_ApplyLimitAndOffset(t *testing.T) {
	articles := domain.ArticleCollection{art1, art2, art3, art4}

	t.Run("complete", func(t *testing.T) {
		assert.Equal(t, articles, articles.ApplyLimitAndOffset(100, 0))
		assert.Equal(t, articles, articles.ApplyLimitAndOffset(4, 0))
		assert.Equal(t, articles, articles.ApplyLimitAndOffset(4, -1))
	})
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, domain.ArticleCollection{}, articles.ApplyLimitAndOffset(100, 10))
		assert.Equal(t, domain.ArticleCollection{}, articles.ApplyLimitAndOffset(3, 4))
		assert.Equal(t, domain.ArticleCollection{}, articles.ApplyLimitAndOffset(-1, 0))
	})
}

func TestArticleHasAuthor(t *testing.T) {
	filter := domain.ArticleHasAuthor("author")
	assert.True(t, filter(domain.Article{Author: domain.User{Name: "author"}}))
	assert.False(t, filter(domain.Article{Author: domain.User{Name: "otherAuthor"}}))
}

func TestArticleHasTag(t *testing.T) {
	filter := domain.ArticleHasTag("tag")
	assert.True(t, filter(domain.Article{TagList: []string{"tag"}}))
	assert.False(t, filter(domain.Article{TagList: []string{"otherTag"}}))
}

func TestArticleIsFavoritedBy(t *testing.T) {
	t.Run("with userName", func(t *testing.T) {
		filter := domain.ArticleIsFavoritedBy("user")
		assert.True(t, filter(domain.Article{FavoritedBy: []domain.User{{Name: "user"}}}))
		assert.False(t, filter(domain.Article{FavoritedBy: []domain.User{{Name: "otherUser"}}}))
		assert.False(t, filter(domain.Article{FavoritedBy: []domain.User{{Name: ""}}}))
	})

	t.Run("without username", func(t *testing.T) {
		emptyFilter := domain.ArticleIsFavoritedBy("")
		assert.False(t, emptyFilter(domain.Article{FavoritedBy: []domain.User{{Name: ""}}})) // always returns false
	})
}
