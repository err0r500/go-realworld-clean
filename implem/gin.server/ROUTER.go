package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"strings"

	"errors"

	"github.com/err0r500/go-realworld-clean/uc"
	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type RouterHandler struct {
	ucHandler   uc.Handler
	authHandler uc.AuthHandler
	Logger      uc.Logger
}

func NewRouter(i uc.Handler, auth uc.AuthHandler) RouterHandler {
	return RouterHandler{
		ucHandler:   i,
		authHandler: auth,
	}
}

func NewRouterWithLogger(i uc.Handler, auth uc.AuthHandler, logger uc.Logger) RouterHandler {
	return RouterHandler{
		ucHandler:   i,
		authHandler: auth,
		Logger:      logger,
	}
}

func (rH RouterHandler) SetRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(rH.errorCatcher())

	rH.profileRoutes(api)
	rH.usersRoutes(api)
	rH.articlesRoutes(api)

	api.GET("/tags", initTrace("http:get_tags"), rH.tagsGet) //tags
}

func (rH RouterHandler) profileRoutes(api *gin.RouterGroup) {
	profiles := api.Group("/profiles")
	profiles.GET("/:username", initTrace("http:get_profile_uname"), rH.profileGet)                                                  // Get a profile of a user of the system. Auth is optional
	profiles.POST("/:username/follow", initTrace("http:post_profile_uname_follow"), rH.jwtMiddleware(), rH.profileFollowPost)       // Follow a user by username
	profiles.DELETE("/:username/follow", initTrace("http:delete_profile_uname_follow"), rH.jwtMiddleware(), rH.profileFollowDelete) // Unfollow a user by username
}

func (rH RouterHandler) usersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	users.POST("", initTrace("http:post_users"), rH.userPost)                  // Register a new user
	users.POST("/login", initTrace("http:post_users_login"), rH.userLoginPost) // Login for existing user

	user := api.Group("/user")
	user.GET("", initTrace("http:get_user"), rH.jwtMiddleware(), rH.userGet)       // Gets the currently logged-in user
	user.PUT("", initTrace("http:put_user"), rH.jwtMiddleware(), rH.userPatch)     // WARNING : it's in fact a PATCH request in the API contract !!!
	user.PATCH("", initTrace("http:patch_user"), rH.jwtMiddleware(), rH.userPatch) // just in case it's fixed one day....
}

func (rH RouterHandler) articlesRoutes(api *gin.RouterGroup) {
	articles := api.Group("/articles")

	articles.GET("", initTrace("http:get_articles"), rH.articlesFilteredGet)
	articles.POST("", initTrace("http:post_articles"), rH.jwtMiddleware(), rH.articlePost)
	articles.GET("/:slug", func(c *gin.Context) { // ugly api contract !
		if c.Param("slug") == "feed" {
			initTrace("http:get_articles_feed")(c)
			rH.jwtMiddleware()(c)
			rH.articlesFeedGet(c)
		} else {
			initTrace("http:get_articles_slug")(c)
			rH.articleGet(c)
		}
	})

	articles.PUT("/:slug", initTrace("http:put_articles_slug"), rH.jwtMiddleware(), rH.articlePut)
	articles.DELETE("/:slug", initTrace("http:delete_articles_slug"), rH.jwtMiddleware(), rH.articleDelete)
	articles.POST("/:slug/favorite", initTrace("http:post_articles_slug_fav"), rH.jwtMiddleware(), rH.updateFavorite)
	articles.DELETE("/:slug/favorite", initTrace("http:delete_articles_slug_fav"), rH.jwtMiddleware(), rH.updateFavorite)

	articles.GET("/:slug/comments", initTrace("http:get_articles_slug_comments"), rH.commentsGet)
	articles.POST("/:slug/comments", initTrace("http:post_articles_slug_comments"), rH.jwtMiddleware(), rH.commentPost)
	articles.DELETE("/:slug/comments/:id", initTrace("http:delete_articles_slug_comments_id"), rH.jwtMiddleware(), rH.commentDelete)
}

const userNameKey = "userNameKey"

func initTrace(name string) gin.HandlerFunc {
	return opengintracing.NewSpan(opentracing.GlobalTracer(), name)
}

func (rH RouterHandler) jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		opengintracing.SpanFromContext(opentracing.GlobalTracer(), "jwt check", false)(c)

		jwt, err := getJWT(c.GetHeader("Authorization"))
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		userName, ok := rH.authHandler.GetUserName(c, jwt)
		if !ok {
			c.Status(http.StatusUnauthorized)
			return
		}

		c.SetAccepted()
		c.Set(userNameKey, userName)
		c.Next()
	}
}

func (rH RouterHandler) getUserNameFromToken(c *gin.Context) string {
	jwt, err := getJWT(c.GetHeader("Authorization"))
	if err != nil {
		return ""
	}

	userName, ok := rH.authHandler.GetUserName(c.Request.Context(), jwt)
	if !ok {
		return ""
	}

	return userName
}

func getJWT(authHeader string) (string, error) {
	splitted := strings.Split(authHeader, "Token ")
	if len(splitted) != 2 {
		return "", errors.New("malformed header")
	}
	return splitted[1], nil
}

func (rH RouterHandler) errorCatcher() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() > 399 {
			c.Render(
				c.Writer.Status(),
				render.Data{
					ContentType: "application/json; charset=utf-8",
					Data:        []byte(`{"errors": {"body": ["wooops, something went wrong !"]}}`),
				},
			)
		}
	}
}

func (RouterHandler) getUserName(c *gin.Context) string {
	if userName, ok := c.Keys[userNameKey].(string); ok {
		return userName
	}
	return ""
}

// log is used to "partially apply" the title to the rH.logger.Log function
// so we can see in the logs from which route the log comes from
func (rH RouterHandler) log(title string) func(...interface{}) {
	return func(logs ...interface{}) {
		rH.Logger.Log(title, logs)
	}
}

func (RouterHandler) MethodAndPath(c *gin.Context) string {
	return fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
}

func startChildSpanFromGinCtx(c *gin.Context, name string) (opentracing.Span, context.Context) {
	parentSpan, ok := opengintracing.GetSpan(c)
	if ok {
		sp := opentracing.StartSpan(
			name,
			opentracing.ChildOf(parentSpan.Context()),
		)
		return sp, opentracing.ContextWithSpan(c, sp)
	}

	sp := opentracing.StartSpan(name)
	return sp, opentracing.ContextWithSpan(c, sp)
}

func logErr(sp opentracing.Span, err error) {
	sp.SetTag("error", true)
	sp.LogFields(log.Error(err))
}

func setStatus(status int, c *gin.Context, sp opentracing.Span) {
	c.Status(status)
	sp.SetTag("http_status", status)
}

func respJSON(status int, body interface{}, c *gin.Context, sp opentracing.Span) {
	sp.SetTag("http_status", status)
	c.JSON(status, body)
}
