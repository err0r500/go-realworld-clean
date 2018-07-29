package server

import (
	"net/http"

	"github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) profileGet(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))
	requestingUserName := ""
	if userName, err := rH.authHandler.GetUserName(c.GetHeader("Authorization")); err == nil {
		requestingUserName = userName
	}

	user, follows, err := rH.ucHandler.ProfileGet(requestingUserName, c.Param("username"))
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": formatter.NewProfileFromDomain(*user, follows)})
}
