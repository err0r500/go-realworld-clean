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
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:user_login")
	defer sp.Finish()

	body := &userLoginPostBody{}
	if err := c.BindJSON(body); err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	user, token, err := rH.ucHandler.UserLogin(ctx, body.User.Email, body.User.Password)
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusOK, gin.H{"user": formatter.NewUserResp(*user, token)}, c, sp)
}
