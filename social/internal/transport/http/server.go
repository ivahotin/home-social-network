package http

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewServer(addr string, endpoints *Endpoints) *http.Server {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("./social/internal/views/html/*.html")

	router.Use(static.Serve("/js", static.LocalFile("./social/internal/views/js", true)))
	router.Use(static.Serve("/css", static.LocalFile("./social/internal/views/css", true)))

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
	}

	router.GET("/", endpoints.Auth.JWTMiddleWare.MiddlewareFunc(), endpoints.Home)
	router.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}