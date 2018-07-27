package server

import (
	"net/http"

	"github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) profileFollowPost(c *gin.Context) {
	log := rH.log(c.Request.URL.Path)

	user, err := rH.ucHandler.ProfileUpdateFollow(rH.getUserName(c), c.Param("username"), true)
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, formatter.NewProfileFromDomain(*user, true))
}
