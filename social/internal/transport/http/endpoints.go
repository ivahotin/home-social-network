package http

import (
	"errors"
	"example.com/social/internal/domain"
	"example.com/social/internal/usecases"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type Endpoints struct {
	Home gin.HandlerFunc
	Auth *AuthEndpoints
}

type AuthEndpoints struct {
	LoginPage gin.HandlerFunc
	RegistrationPage gin.HandlerFunc
	JWTMiddleWare *jwt.GinJWTMiddleware
	SignUp gin.HandlerFunc
}

func MakeEndpoints(
	signUpUseCase usecases.SignUpUseCase,
	signInUseCase usecases.SignInUseCase) *Endpoints {
	return &Endpoints{
		Home: makeHomePage(),
		Auth: &AuthEndpoints{
			LoginPage: makeLoginPage(),
			RegistrationPage: makeRegistrationPage(),
			JWTMiddleWare: getJwtMiddleware(signInUseCase),
			SignUp: makeSignUpEndpoint(signUpUseCase),
		},
	}
}

func makeHomePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}

func makeLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	}
}

func makeRegistrationPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	}
}

func makeSignUpEndpoint(signUseCase usecases.SignUpUseCase) gin.HandlerFunc {
	var signUpRequest SignUpRequest
	return func(ctx *gin.Context) {
		if err := ctx.ShouldBind(&signUpRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		credentials, err := domain.NewCredentials(signUpRequest.Username, signUpRequest.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Try to repeat later"})
			return
		}

		var gender domain.GenderType
		switch signUpRequest.Gender {
		case "male": gender = domain.Male
		case "female": gender = domain.Female
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid gender"})
			return
		}
		profile := domain.Profile{
			Username: credentials.Username,
			Password: credentials.Password,
			Firstname: signUpRequest.Firstname,
			Lastname: signUpRequest.Lastname,
			Age: signUpRequest.Age,
			Gender: gender,
			Interests: signUpRequest.Interests,
			City: signUpRequest.City,
		}

		err = signUseCase.SignUp(&profile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something goes wrong"})
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func makeSignInEndpoint(signInUseCase usecases.SignInUseCase) func(*gin.Context) (interface{}, error) {
	var signInRequest SignInRequest
	return func (ctx *gin.Context) (interface{}, error) {
		if err := ctx.ShouldBind(&signInRequest); err != nil {
			return "", jwt.ErrMissingLoginValues
		}

		credentials, err := domain.NewCredentials(signInRequest.Username, signInRequest.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Try to repeat later"})
			return nil, err
		}

		isMatch, err := signInUseCase.SignIn(credentials)
		switch {
		case errors.Is(err, domain.ProfileNotFound):
			return nil, jwt.ErrFailedAuthentication
		case err != nil:
			return nil, err
		}

		if !isMatch {
			return nil, jwt.ErrFailedAuthentication
		}

		return &User{ UserName: credentials.Username }, nil
	}
}

func getJwtMiddleware(signInUseCase usecases.SignInUseCase) *jwt.GinJWTMiddleware {
	const identityKey = "id"
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("84636aa7-1b02-47b7-8993-18b1598d8408"),
		Timeout:     time.Minute,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: makeSignInEndpoint(signInUseCase),
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.Redirect(http.StatusFound, "/auth/sign-in")
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,

		SendCookie:       true,
		SecureCookie: 	  false,
		CookieHTTPOnly:   true,
		CookieSameSite:   http.SameSiteDefaultMode,
	})

	if err != nil {
		log.Fatal("Jwt initialization error: ", err.Error())
	}

	errInit := middleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	return middleware
}
