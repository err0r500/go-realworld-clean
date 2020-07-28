package server

import (
	"net/http"

	"github.com/err0r500/go-realworld-clean/domain"

	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"

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
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:post_user")
	defer sp.Finish()

	body := &userPostRequest{}
	if err := c.BindJSON(body); err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	user, token, err := rH.ucHandler.UserCreate(ctx, body.User.Username, body.User.Email, body.User.Password)
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusCreated, gin.H{"user": formatter.NewUserResp(*user, token)}, c, sp)
}

func (rH RouterHandler) userGet(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:get_current_user")
	defer sp.Finish()

	user, token, err := rH.ucHandler.UserGet(ctx, rH.getUserName(c))
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusOK, gin.H{"user": formatter.NewUserResp(*user, token)}, c, sp)
}

// New user details.
type userPutRequest struct {
	User struct {
		Email    *string `json:"email,omitempty"`
		Name     *string `json:"username,omitempty"`
		Bio      *string `json:"bio,omitempty"`
		Image    *string `json:"image,omitempty"`
		Password *string `json:"password,omitempty"`
	} `json:"user,required"`
}

func (req userPutRequest) getEditableFields() map[domain.UserUpdatableProperty]*string {
	return map[domain.UserUpdatableProperty]*string{
		domain.UserEmail:     req.User.Email,
		domain.UserName:      req.User.Name,
		domain.UserBio:       req.User.Bio,
		domain.UserImageLink: req.User.Image,
		domain.UserPassword:  req.User.Password,
	}
}

func (rH RouterHandler) userPatch(c *gin.Context) {
	sp, ctx := startChildSpanFromGinCtx(c, "http_handler:patch_user")
	defer sp.Finish()

	req := &userPutRequest{}
	if err := c.BindJSON(req); err != nil {
		logErr(sp, err)
		setStatus(http.StatusBadRequest, c, sp)
		return
	}

	user, token, err := rH.ucHandler.UserEdit(ctx, rH.getUserName(c), req.getEditableFields())
	if err != nil {
		logErr(sp, err)
		setStatus(http.StatusUnprocessableEntity, c, sp)
		return
	}

	respJSON(http.StatusOK, gin.H{"user": formatter.NewUserResp(*user, token)}, c, sp)
}
