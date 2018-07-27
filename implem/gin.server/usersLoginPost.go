package server

import (
	"net/http"

	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

type userLoginPostBody struct {
	User struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	} `json:"user" binding:"required"`
}

func (rH RouterHandler) userLoginPost(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	body := &userLoginPostBody{}
	if err := c.BindJSON(body); err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	user, token, err := rH.ucHandler.UserLogin(body.User.Email, body.User.Password)
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": formatter.NewUserResp(*user, token)})
}
