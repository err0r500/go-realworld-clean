package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (rH RouterHandler) tagsGet(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	tags, err := rH.ucHandler.Tags()
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	if tags == nil {
		tags = []string{}
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}
