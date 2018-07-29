package server

import (
	"net/http"

	"github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) profileFollowDelete(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	user, err := rH.ucHandler.ProfileUpdateFollow(rH.getUserName(c), c.Param("username"), false)
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": formatter.NewProfileFromDomain(*user, false)})
}
