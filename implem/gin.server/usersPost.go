package server

import (
	"net/http"

	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/opentracing/opentracing-go"

	"github.com/gin-gonic/gin"
)

type userPostRequest struct {
	User struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	} `json:"user" binding:"required"`
}

func (rH RouterHandler) userPost(c *gin.Context) {
	log := rH.log(rH.MethodAndPath(c))

	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "http:post_user")
	defer span.Finish()

	body := &userPostRequest{}
	if err := c.BindJSON(body); err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	user, token, err := rH.ucHandler.UserCreate(ctx, body.User.Username, body.User.Email, body.User.Password)
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": formatter.NewUserResp(*user, token)})
}
