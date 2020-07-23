package server

import (
	"net/http"

	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/opentracing/opentracing-go"

	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) profileGet(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "http:get_profile")
	defer span.Finish()

	user, follows, err := rH.ucHandler.ProfileGet(ctx, rH.getUserNameFromToken(c), c.Param("username"))
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": formatter.NewProfileFromDomain(*user, follows)})
}
