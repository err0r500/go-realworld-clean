package server

import (
	"net/http"

	"strconv"

	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/err0r500/go-realworld-clean/uc"
	"github.com/gin-gonic/gin"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
)

func (rH RouterHandler) articlesFilteredGet(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:get_articles_filtered")
	defer sp.Finish()

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = defaultOffset
	}

	user, articles, count, err := rH.ucHandler.GetArticles(ctx,
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
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusOK, gin.H{"articles": formatter.NewArticlesFromDomain(user, articles...), "articlesCount": count}, c, sp)
}

func (rH RouterHandler) articlesFeedGet(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:get_articles_feed")
	defer sp.Finish()

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = defaultOffset
	}

	user, articles, count, err := rH.ucHandler.ArticlesFeed(ctx, rH.getUserNameFromToken(c), limit, offset)
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusOK, gin.H{"articles": formatter.NewArticlesFromDomain(user, articles...), "articlesCount": count}, c, sp)
}
