package server

import (
	"net/http"

	"github.com/opentracing/opentracing-go"

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

	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "http:post_user_login")
	defer span.Finish()

	body := &userLoginPostBody{}
	if err := c.BindJSON(body); err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	user, token, err := rH.ucHandler.UserLogin(ctx, body.User.Email, body.User.Password)
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": formatter.NewUserResp(*user, token)})
}
