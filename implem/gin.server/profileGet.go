package server

import (
	"net/http"

	"github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) profileGet(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	profile, err := rH.ucHandler.ProfileGet(c.Param("username"))
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": formatter.NewProfileFromDomain(*profile)})
}
