package server

import (
	"net/http"

	"strconv"

	"github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) commentsGet(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	comments, err := rH.ucHandler.CommentsGet(c.Param("slug"))
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": formatter.NewCommentsFromDomain(comments...)})
}

type commentPostReq struct {
	Comment struct {
		Body string `json:"body,required"`
	} `json:"comment,required"`
}

func (rH RouterHandler) commentPost(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	req := &commentPostReq{}
	if err := c.BindJSON(req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	comment, err := rH.ucHandler.CommentsPost(rH.getUserName(c), c.Param("slug"), req.Comment.Body)
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"comment": formatter.NewCommentFromDomain(*comment)})
}

func (rH RouterHandler) commentDelete(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := rH.ucHandler.CommentsDelete(rH.getUserName(c), c.Param("slug"), id); err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.Status(http.StatusOK)
}
