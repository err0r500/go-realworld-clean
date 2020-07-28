package server

import (
	"net/http"

	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) updateFavorite(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:update_favorites")
	defer sp.Finish()

	favorite := true
	switch c.Request.Method {
	case "POST":
		break
	case "DELETE":
		favorite = false
	default: // should not be testable :) for regression only
		setStatus(http.StatusBadRequest, c, sp)
		return
	}

	user, article, err := rH.ucHandler.FavoritesUpdate(ctx, rH.getUserName(c), c.Param("slug"), favorite)
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusOK, gin.H{"article": formatter.NewArticleFromDomain(*article, user)}, c, sp)
}
