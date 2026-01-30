package services

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"banking/internal/database"
	"banking/internal/models"
	"banking/internal/repositories"
	"banking/internal/utils"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	db       database.Service
	userRepo *repositories.UserRepository
}

func NewAuthService(db database.Service) *AuthService {
	userRepo := repositories.NewUserRepository(db.DB())
	return &AuthService{
		db:       db,
		userRepo: userRepo,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user1@test.com"`
	Password string `json:"password" binding:"required" example:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"DcwAT..."`
}

type UserResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string    `json:"email" example:"user1@test.com"`
	CreatedAt time.Time `json:"created_at" example:"2026-01-30T12:00:00Z"`
}

type TokenPairResponse struct {
	Code    int    `json:"code" example:"200"`
	Expire  string `json:"expire" example:"2026-01-31T12:00:00Z"`
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Message string `json:"message" example:"success"`
}

// MeHandler godoc
// @Summary Get current user
// @Description Returns current authenticated user information
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} services.UserResponse
// @Failure 401 {object} services.ErrorResponse
// @Router /api/auth/me [get]
// @Security BearerAuth
func (s *AuthService) MeHandler(c *gin.Context) {
	user := utils.MustGetUserFromContext(c)

	response := UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// Middleware
func (s *AuthService) JWTInitParams() *jwt.GinJWTMiddleware {
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is required")
	}

	return &jwt.GinJWTMiddleware{
		Realm:      "test",
		Key:        []byte(jwtSecret),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,

		// Add 60 seconds leeway for clock skew tolerance
		ParseOptions: []gojwt.ParserOption{
			gojwt.WithLeeway(60 * time.Second),
		},

		IdentityKey: utils.IdentityKey,

		RefreshTokenCookieName: "refresh_token",
		CookieName:             "token",
		SendCookie:             true,
		CookieHTTPOnly:         true,
		CookieMaxAge:           time.Hour * 24,

		PayloadFunc:     s.createPayloadFunc(),
		IdentityHandler: s.createIdentityHandler(),
		Authenticator:   s.createAuthenticator(),
		Authorizer:      s.createAuthorizator(),
		Unauthorized:    s.createUnauthorized(),
		LogoutResponse:  s.createLogoutResponse(),
		TokenLookup:     "header: Authorization, cookie: token",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	}
}

func (s *AuthService) createPayloadFunc() func(data any) gojwt.MapClaims {
	return func(data any) gojwt.MapClaims {
		if v, ok := data.(*models.User); ok {
			return gojwt.MapClaims{
				utils.IdentityKey: v.ID.String(),
			}
		}
		return gojwt.MapClaims{}
	}
}

func (s *AuthService) createIdentityHandler() func(c *gin.Context) any {
	return func(c *gin.Context) any {
		claims := jwt.ExtractClaims(c)
		if len(claims) == 0 {
			return nil
		}

		userIDStr, ok := claims[utils.IdentityKey].(string)
		if !ok || userIDStr == "" {
			return nil
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil
		}

		ctx := c.Request.Context()
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil || user == nil {
			// trigger 401
			return nil
		}

		return user
	}
}

func (s *AuthService) createAuthenticator() func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		var loginVals LoginRequest
		if err := c.ShouldBindJSON(&loginVals); err != nil {
			return nil, fmt.Errorf("invalid login data: %s", utils.FormatValidationErrors(err))
		}

		ctx := c.Request.Context()
		user, err := s.userRepo.GetByEmail(ctx, loginVals.Email)
		if err != nil {
			if err == repositories.ErrUserNotFound {
				return nil, jwt.ErrFailedAuthentication
			}
			return nil, fmt.Errorf("database error: %w", err)
		}

		if !utils.CheckPasswordHash(loginVals.Password, user.PasswordHash) {
			return nil, jwt.ErrFailedAuthentication
		}

		return user, nil
	}
}

func (s *AuthService) createAuthorizator() func(c *gin.Context, data any) bool {
	return func(c *gin.Context, data any) bool {
		// We dont have roles
		return true
	}
}

func (s *AuthService) createUnauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		// Ensure 401 status code for unauthorized access
		if code == 0 {
			code = http.StatusUnauthorized
		}
		if message == "" {
			message = "Unauthorized"
		}
		utils.RespondWithError(c, code, message, nil)
		c.Abort()
	}
}

func (s *AuthService) createLogoutResponse() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		user, exists := c.Get(utils.IdentityKey)

		response := gin.H{
			"code":    http.StatusOK,
			"message": "Successfully logged out",
		}

		if len(claims) > 0 {
			response["logged_out_user"] = claims[utils.IdentityKey]
		}
		if exists {
			if userModel, ok := user.(*models.User); ok {
				response["user_info"] = userModel.Email
			}
		}

		c.JSON(http.StatusOK, response)
	}
}
