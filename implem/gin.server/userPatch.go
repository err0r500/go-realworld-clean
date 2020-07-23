package server

import (
	"net/http"

	"github.com/opentracing/opentracing-go"

	"github.com/err0r500/go-realworld-clean/domain"
	formatter "github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

//New user details.
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
	log := rH.log(rH.MethodAndPath(c))

	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "http:patch_user")
	defer span.Finish()

	req := &userPutRequest{}
	if err := c.BindJSON(req); err != nil {
		log(err)
		c.Status(http.StatusBadRequest)
		return
	}

	user, token, err := rH.ucHandler.UserEdit(ctx, rH.getUserName(c), req.getEditableFields())
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": formatter.NewUserResp(*user, token)})
}
