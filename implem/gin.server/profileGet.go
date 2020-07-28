package server

import (
	"net/http"

	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) profileGet(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:profile_get")
	defer sp.Finish()

	user, follows, err := rH.ucHandler.ProfileGet(ctx, rH.getUserNameFromToken(c), c.Param("username"))
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusOK, gin.H{"profile": formatter.NewProfileFromDomain(*user, follows)}, c, sp)
}
