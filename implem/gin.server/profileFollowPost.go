package server

import (
	"net/http"

	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/opentracing/opentracing-go"

	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) profileFollowPost(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "http:post_profile_follow")
	defer span.Finish()

	user, err := rH.ucHandler.ProfileUpdateFollow(ctx, rH.getUserName(c), c.Param("username"), true)
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": formatter.NewProfileFromDomain(*user, false)})
}
