package server

import (
	"net/http"

	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/opentracing/opentracing-go"

	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) userGet(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "http:get_user")
	defer span.Finish()

	user, token, err := rH.ucHandler.UserGet(ctx, rH.getUserName(c))
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": formatter.NewUserResp(*user, token)})
}
