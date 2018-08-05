package server

import (
	"net/http"

	"strconv"

	"github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/err0r500/go-realworld-clean/uc"
	"github.com/gin-gonic/gin"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
)

func (rH RouterHandler) articlesFilteredGet(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = defaultOffset
	}

	user, articles, count, err := rH.ucHandler.GetArticles(
		rH.getUserNameFromToken(c),
		limit,
		offset,
		uc.NewFilters(
			c.Query("author"),
			c.Query("tag"),
			c.Query("favorited"),
		),
	)
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": formatter.NewArticlesFromDomain(user, articles...), "articlesCount": count})
}

func (rH RouterHandler) articlesFeedGet(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = defaultOffset
	}

	user, articles, count, err := rH.ucHandler.ArticlesFeed(rH.getUserNameFromToken(c), limit, offset)
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": formatter.NewArticlesFromDomain(user, articles...), "articlesCount": count})
}
