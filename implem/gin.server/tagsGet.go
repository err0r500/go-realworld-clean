package server

import (
	"net/http"

	"github.com/opentracing/opentracing-go"

	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) tagsGet(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "http:get_tag")
	defer span.Finish()

	tags, err := rH.ucHandler.Tags(ctx)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	if tags == nil {
		tags = []string{}
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}
