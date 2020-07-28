package server

import (
	"net/http"

	"github.com/err0r500/go-realworld-clean/domain"
	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

type ArticleReq struct {
	Article struct {
		Title       *string   `json:"title"`
		Description *string   `json:"description"`
		Body        *string   `json:"body"`
		TagList     *[]string `json:"tagList"`
	} `json:"article,required"`
}

func (req ArticleReq) getEditableFields() map[domain.ArticleUpdatableField]*string {
	return map[domain.ArticleUpdatableField]*string{
		domain.Title:       req.Article.Title,
		domain.Description: req.Article.Description,
		domain.Body:        req.Article.Body,
	}
}

func articleFromReq(req *ArticleReq) domain.Article {
	return domain.Article{
		Title:       *req.Article.Title,
		Description: *req.Article.Description,
		Body:        *req.Article.Body,
		TagList:     *req.Article.TagList,
	}
}

func (rH RouterHandler) articlePost(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:post_articles")
	defer sp.Finish()

	req := &ArticleReq{}
	if err := c.BindJSON(req); err != nil {
		logErr(sp, err)
		setStatus(http.StatusBadRequest, c, sp)
		return
	}

	user, article, err := rH.ucHandler.ArticlePost(ctx, rH.getUserName(c), articleFromReq(req))
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusCreated, gin.H{"article": formatter.NewArticleFromDomain(*article, user)}, c, sp)
}

func (rH RouterHandler) articlePut(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:put_articles")
	defer sp.Finish()

	req := &ArticleReq{}
	if err := c.BindJSON(req); err != nil {
		logErr(sp, err)
		setStatus(http.StatusBadRequest, c, sp)
		return
	}

	user, article, err := rH.ucHandler.ArticlePut(ctx, rH.getUserName(c), c.Param("slug"), req.getEditableFields())
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusOK, gin.H{"article": formatter.NewArticleFromDomain(*article, user)}, c, sp)
}

func (rH RouterHandler) articleGet(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:get_articles")
	defer sp.Finish()

	user, article, err := rH.ucHandler.ArticleGet(ctx, rH.getUserName(c), c.Param("slug"))
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}
	if article == nil {
		setStatus(http.StatusNotFound, c, sp)
		return
	}

	respJSON(http.StatusOK, gin.H{"article": formatter.NewArticleFromDomain(*article, user)}, c, sp)
}

func (rH RouterHandler) articleDelete(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:delete_articles")
	defer sp.Finish()

	if err := rH.ucHandler.ArticleDelete(ctx, rH.getUserName(c), c.Param("slug")); err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	setStatus(http.StatusOK, c, sp)
}
