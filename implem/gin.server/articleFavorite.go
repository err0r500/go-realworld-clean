package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) updateFavorite(c *gin.Context) {
	log := rH.log(c.Request.URL.Path)

	favorite := true
	switch c.Request.Method {
	case "POST":
		break
	case "DELETE":
		favorite = false
	default: // should not be testable :) for regression only
		c.Status(http.StatusBadRequest)
		return
	}

	userName, err := rH.getUserName(c)
	if err != nil {
		log(err)
		c.Status(http.StatusUnauthorized)
		return
	}

	article, err := rH.ucHandler.FavoritesUpdate(userName, c.Param("slug"), favorite)
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}
