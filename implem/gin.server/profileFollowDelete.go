package server

import (
	"net/http"

	"github.com/opentracing/opentracing-go"

	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) profileFollowDelete(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "http:delete_profile_follows")
	defer span.Finish()

	user, err := rH.ucHandler.ProfileUpdateFollow(ctx, rH.getUserName(c), c.Param("username"), false)
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": formatter.NewProfileFromDomain(*user, false)})
}
