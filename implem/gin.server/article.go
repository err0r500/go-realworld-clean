package server

import (
	"net/http"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/implem/json.formatter"
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
	log := rH.log(rH.MethodAndPath(c))

	req := &ArticleReq{}
	if err := c.BindJSON(req); err != nil {
		log(err)
		c.Status(http.StatusBadRequest)
		return
	}

	user, article, err := rH.ucHandler.ArticlePost(rH.getUserName(c), articleFromReq(req))
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"article": formatter.NewArticleFromDomain(*article, user)})
}

func (rH RouterHandler) articlePut(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	req := &ArticleReq{}
	if err := c.BindJSON(req); err != nil {
		log(err)
		c.Status(http.StatusBadRequest)
		return
	}
	user, article, err := rH.ucHandler.ArticlePut(rH.getUserName(c), c.Param("slug"), req.getEditableFields())
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": formatter.NewArticleFromDomain(*article, user)})
}

func (rH RouterHandler) articleGet(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	user, article, err := rH.ucHandler.ArticleGet(rH.getUserName(c), c.Param("slug"))
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	c.JSON(http.StatusOK, gin.H{"article": formatter.NewArticleFromDomain(*article, user)})
}

func (rH RouterHandler) articleDelete(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	if err := rH.ucHandler.ArticleDelete(rH.getUserName(c), c.Param("slug")); err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	c.Status(http.StatusOK)
}
