package services

import (
	"net/http"
	"time"

	"banking/internal/database"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	db database.Service
}

func NewAuthService(db database.Service) *AuthService {
	return &AuthService{db: db}
}

// Types
const (
	identityKey = "id"
	userAdmin   = "admin"
)

type User struct {
	UserName  string
	FirstName string
	LastName  string
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// LoginHandler godoc
// @Summary Login
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body login true "Login credentials"
// @Success 200 {object} map[string]string
// @Router /api/auth/login [post]
func (s *AuthService) LoginHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Login test"

	c.JSON(http.StatusOK, resp)
}

// MeHandler godoc
// @Summary Get current user
// @Description Returns current authenticated user information
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/auth/me [get]
// @Security BearerAuth
func (s *AuthService) MeHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Me test"

	c.JSON(http.StatusOK, resp)
}

// Middleware
func (s *AuthService) JWTInitParams() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "identityKey",
		PayloadFunc: payloadFunc(),

		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(),
		Authorizer:      authorizator(),
		Unauthorized:    unauthorized(),
		LogoutResponse:  logoutResponse(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}

func payloadFunc() func(data any) gojwt.MapClaims {
	return func(data any) gojwt.MapClaims {
		if v, ok := data.(*User); ok {
			return gojwt.MapClaims{
				identityKey: v.UserName,
			}
		}
		return gojwt.MapClaims{}
	}
}

func identityHandler() func(c *gin.Context) any {
	return func(c *gin.Context) any {
		claims := jwt.ExtractClaims(c)
		return &User{
			UserName: claims[identityKey].(string),
		}
	}
}

func authenticator() func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		var loginVals login
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		userID := loginVals.Username
		password := loginVals.Password

		if (userID == userAdmin && password == userAdmin) ||
			(userID == "test" && password == "test") {
			return &User{
				UserName:  userID,
				LastName:  "Bo-Yi",
				FirstName: "Wu",
			}, nil
		}
		return nil, jwt.ErrFailedAuthentication
	}
}

func authorizator() func(c *gin.Context, data any) bool {
	return func(c *gin.Context, data any) bool {
		if v, ok := data.(*User); ok && v.UserName == "admin" {
			return true
		}
		return false
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func logoutResponse() func(c *gin.Context) {
	return func(c *gin.Context) {
		// This demonstrates that claims are now accessible during logout
		claims := jwt.ExtractClaims(c)
		user, exists := c.Get(identityKey)

		response := gin.H{
			"code":    http.StatusOK,
			"message": "Successfully logged out",
		}

		// Show that we can access user information during logout
		if len(claims) > 0 {
			response["logged_out_user"] = claims[identityKey]
		}
		if exists {
			response["user_info"] = user.(*User).UserName
		}

		c.JSON(http.StatusOK, response)
	}
}
