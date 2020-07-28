package server

import (
	"net/http"

	"strconv"

	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) commentsGet(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:get_comments")
	defer sp.Finish()

	comments, err := rH.ucHandler.CommentsGet(ctx, c.Param("slug"))
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusOK, gin.H{"comments": formatter.NewCommentsFromDomain(comments...)}, c, sp)
}

type commentPostReq struct {
	Comment struct {
		Body string `json:"body,required"`
	} `json:"comment,required"`
}

func (rH RouterHandler) commentPost(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:post_comments")
	defer sp.Finish()

	req := &commentPostReq{}
	if err := c.BindJSON(req); err != nil {
		logErr(sp, err)
		setStatus(http.StatusBadRequest, c, sp)
		return
	}

	comment, err := rH.ucHandler.CommentsPost(ctx, rH.getUserName(c), c.Param("slug"), req.Comment.Body)
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusCreated, gin.H{"comment": formatter.NewCommentFromDomain(*comment)}, c, sp)
}

func (rH RouterHandler) commentDelete(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:delete_comments")
	defer sp.Finish()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusBadRequest, c, sp)
		return
	}

	if err := rH.ucHandler.CommentsDelete(ctx, rH.getUserName(c), c.Param("slug"), id); err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	setStatus(http.StatusOK, c, sp)
}
