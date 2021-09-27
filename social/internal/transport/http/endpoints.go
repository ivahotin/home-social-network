package http

import (
	"errors"
	"example.com/social/internal/domain"
	"example.com/social/internal/usecases"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Endpoints struct {
	Home 	gin.HandlerFunc
	Auth 	*AuthEndpoints
	Profile *ProfileEndpoints
}

type AuthEndpoints struct {
	LoginPage 			gin.HandlerFunc
	RegistrationPage 	gin.HandlerFunc
	JWTMiddleWare 		*jwt.GinJWTMiddleware
	SignUp 				gin.HandlerFunc
}

type ProfileEndpoints struct {
	Me 				gin.HandlerFunc
	SearchProfile 	gin.HandlerFunc
	EditProfile 	gin.HandlerFunc
}

func MakeEndpoints(
	signUpUseCase usecases.SignUpUseCase,
	signInUseCase usecases.SignInUseCase,
	getProfileByUsername usecases.GetProfileGetUsernameUseCase,
	getProfilesBySearchTerm usecases.GetProfilesBySearchTerm) *Endpoints {
	return &Endpoints{
		Home: makeHomePage(),
		Auth: &AuthEndpoints{
			LoginPage: makeLoginPage(),
			RegistrationPage: makeRegistrationPage(),
			JWTMiddleWare: getJwtMiddleware(signInUseCase),
			SignUp: makeSignUpEndpoint(signUpUseCase),
		},
		Profile: &ProfileEndpoints{
			Me: makeMyProfileEndpoint(getProfileByUsername),
			SearchProfile: makeSearchEndpoint(getProfilesBySearchTerm),
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
		switch {
		case errors.Is(err, domain.SuchUsernameExists):
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		case err != nil:
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

		signInResult, err := signInUseCase.SignIn(credentials)
		switch {
		case errors.Is(err, domain.ProfileNotFound):
			return nil, jwt.ErrFailedAuthentication
		case err != nil:
			return nil, err
		}

		if !signInResult.IsMatch {
			return nil, jwt.ErrFailedAuthentication
		}

		ctx.Header("x-user-id", strconv.FormatInt(signInResult.Id, 10))
		return &User{ UserName: credentials.Username, Id: signInResult.Id }, nil
	}
}

func makeMyProfileEndpoint(getProfileByUsernameUseCase usecases.GetProfileGetUsernameUseCase) gin.HandlerFunc {
	return func (ctx *gin.Context) {
		claims := jwt.ExtractClaims(ctx)
		username := claims["username"].(string)
		profile, err := getProfileByUsernameUseCase.GetProfileByUsername(username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		if profile == nil {
			ctx.Status(http.StatusNotFound)
			return
		}

		var gender string
		switch profile.Gender {
		case domain.Male: gender = "Male"
		case domain.Female: gender = "Female"
		}
		ctx.JSON(http.StatusOK, gin.H{
			"firstname": profile.Firstname,
			"lastname": profile.Lastname,
			"age": profile.Age,
			"gender": gender,
			"city": profile.City,
			"interests": profile.Interests,
		})
	}
}

func makeSearchEndpoint(getProfilesBySearchTerm usecases.GetProfilesBySearchTerm) gin.HandlerFunc {
	var searchRequest SearchUsersRequest
	return func (ctx *gin.Context) {
		if err := ctx.ShouldBind(&searchRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		claims := jwt.ExtractClaims(ctx)
		myId := int64(claims["id"].(float64))

		profilesSearchResult, err := getProfilesBySearchTerm.GetProfilesBySearchTerm(
			searchRequest.SearchTerm,
			searchRequest.Cursor,
			searchRequest.Limit,
			myId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		profiles := make([]*Profile, 0, 32)
		for _, domainProfile := range profilesSearchResult.Profiles {
			profile := new(Profile)
			profile.Id 			= domainProfile.Id
			profile.Username 	= domainProfile.Username
			profile.Firstname 	= domainProfile.Firstname
			profile.Lastname 	= domainProfile.Lastname
			profile.Age 		= domainProfile.Age
			profile.City        = domainProfile.City
			switch domainProfile.Gender {
			case domain.Male: profile.Gender = "Male"
			case domain.Female: profile.Gender = "Female"
			}
			profile.Interests 	= domainProfile.Interests

			profiles = append(profiles, profile)
		}

		ctx.JSON(http.StatusOK, GetProfilesBySearchTerm{
			Profiles: 	profiles,
			PrevCursor: profilesSearchResult.PrevCursor,
			NextCursor: profilesSearchResult.NextCursor,
		})
	}
}

func getJwtMiddleware(signInUseCase usecases.SignInUseCase) *jwt.GinJWTMiddleware {
	const identityKey = "id"
	const username = "username"
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("84636aa7-1b02-47b7-8993-18b1598d8408"),
		Timeout:     time.Minute * 5,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					username: v.UserName,
					identityKey: v.Id,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[username].(string),
				Id: int64(claims[identityKey].(float64)),
			}
		},
		Authenticator: makeSignInEndpoint(signInUseCase),
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			if c.FullPath() == "/auth/sign-in" {
				c.JSON(http.StatusUnauthorized, gin.H{"message": message})
				return
			}
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
