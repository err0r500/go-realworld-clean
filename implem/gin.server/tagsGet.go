package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) tagsGet(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:get_tags")
	defer sp.Finish()

	tags, err := rH.ucHandler.Tags(ctx)
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	if tags == nil {
		tags = []string{}
	}
	respJSON(http.StatusOK, gin.H{"tags": tags}, c, sp)
}
