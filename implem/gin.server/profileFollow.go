package server

import (
	"net/http"

	"github.com/opentracing/opentracing-go"

	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) profileFollowPost(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:post_profile_follow")
	defer sp.Finish()

	span, ctx := opentracing.StartSpanFromContext(ctx, "http:post_profile_follow")
	defer span.Finish()

	user, err := rH.ucHandler.ProfileUpdateFollow(ctx, rH.getUserName(c), c.Param("username"), true)
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusOK, gin.H{"profile": formatter.NewProfileFromDomain(*user, false)}, c, sp)
}

func (rH RouterHandler) profileFollowDelete(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:delete_profile_follow")
	defer sp.Finish()

	user, err := rH.ucHandler.ProfileUpdateFollow(ctx, rH.getUserName(c), c.Param("username"), false)
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusOK, gin.H{"profile": formatter.NewProfileFromDomain(*user, false)}, c, sp)
}
