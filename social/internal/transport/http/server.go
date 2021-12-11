package http

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewServer(addr string, endpoints *Endpoints) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("./views/html/*.html")

	router.Use(static.Serve("/js", static.LocalFile("./views/js", true)))
	router.Use(static.Serve("/css", static.LocalFile("./views/css", true)))
	router.Use(
		static.Serve(
			"/favicon.ico",
			static.LocalFile("./views/images/favicon.ico", true),
			),
		)

	auth := router.Group("/auth")
	{
		auth.GET("/sign-in", endpoints.Auth.LoginPage)
		auth.POST("/sign-in", endpoints.Auth.JWTMiddleWare.LoginHandler)
		auth.GET("/sign-up", endpoints.Auth.RegistrationPage)
		auth.POST("sign-up", endpoints.Auth.SignUp)
		auth.POST("/sign-out", endpoints.Auth.JWTMiddleWare.LogoutHandler)
		auth.GET("/refresh-token", endpoints.Auth.JWTMiddleWare.RefreshHandler)
	}

	profile := router.Group("/profiles")
	profile.Use(endpoints.Auth.JWTMiddleWare.MiddlewareFunc())
	{
		profile.GET("/me", endpoints.Profile.Me)
		profile.GET("/:id", endpoints.Profile.GetProfile)
		profile.GET("", endpoints.Profile.SearchProfile)
	}

	following := router.Group("/following")
	following.Use(endpoints.Auth.JWTMiddleWare.MiddlewareFunc())
	{
		following.POST("/:followed/follow", endpoints.Following.Follow)
		following.POST("/:followed/unfollow", endpoints.Following.UnFollow)
	}

	router.GET("/", endpoints.Auth.JWTMiddleWare.MiddlewareFunc(), endpoints.Home)
	router.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	chat := router.Group("/chats")
	chat.Use(endpoints.Auth.JWTMiddleWare.MiddlewareFunc())
	{
		chat.POST("/:chat_id/messages", endpoints.Chat.Publish)
	}

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}