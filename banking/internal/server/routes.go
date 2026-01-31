package server

import (
	"log"
	"net/http"

	"banking/internal/services"
	"banking/internal/validators"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := validators.RegisterCustomValidators(v); err != nil {
			log.Fatal("Failed to register custom validators: " + err.Error())
		}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://banking-demo.gornostay25.dev"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/health", s.healthHandler)

	public := r.Group("/api")

	// AUTH
	authService := services.NewAuthService(s.db)
	authMiddleware, err := jwt.New(authService.JWTInitParams())
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	authGroup := public.Group("/auth")
	{
		authGroup.POST("/login", s.loginHandler(authMiddleware))
		authGroup.POST("/refresh", s.refreshHandler(authMiddleware))
		authGroup.POST("/logout", s.logoutHandler(authMiddleware))
		authGroup.GET("/me", authMiddleware.MiddlewareFunc(), authService.MeHandler)
	}

	protectedGroup := public.Group("")
	protectedGroup.Use(authMiddleware.MiddlewareFunc())

	// Account Operations
	accountService := services.NewAccountService(s.db)
	accountGroup := protectedGroup.Group("/accounts")
	{
		accountGroup.GET("", accountService.ListAccountsHandler)
		accountGroup.GET("/:id/balance", accountService.GetAccountBalanceHandler)
	}

	// Transaction Operations
	transactionService := services.NewTransactionService(s.db)
	transactionGroup := protectedGroup.Group("/transactions")
	{
		transactionGroup.POST("/transfer", transactionService.TransferHandler)
		transactionGroup.POST("/exchange", transactionService.ExchangeHandler)
		transactionGroup.GET("", transactionService.GetTransactionsHandler)
	}

	// Swagger documentation
	r.GET("/swagger/*any", func(ctx *gin.Context) {
		// Redirect to swagger index.html
		if ctx.Request.URL.Path == "/swagger" || ctx.Request.URL.Path == "/swagger/" {
			ctx.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
			return
		}
		ctx.Next()
	}, ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

// loginHandler godoc
// @Summary Login
// @Description Authenticate user and get JWT token. Returns authentication tokens as HTTP-only cookies: `token` and `refresh_token`.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body services.LoginRequest true "Login credentials"
// @Success 200 {object} services.TokenPairResponse
// @header 200 {string} Set-Cookie "Sets two HTTP-only cookies: token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9... and refresh_token=NGDcrC2V..."
// @Failure 401 {object} services.ErrorResponse
// @Router /api/auth/login [post]
func (s *Server) loginHandler(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return authMiddleware.LoginHandler
}

// refreshHandler godoc
// @Summary Refresh token
// @Description Refresh JWT access token using refresh token. Returns new access token and refresh token. Also sets HTTP-only cookies: `token` and `refresh_token`.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} services.TokenPairResponse
// @header 200 {string} Set-Cookie "Sets two HTTP-only cookies: token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9... and refresh_token=Jk1WWhw-eo6NWpDH5w3p4ky69gjrCZeuRfG5_2rFsvE="
// @Failure 401 {object} services.ErrorResponse
// @Router /api/auth/refresh [post]
// @Security BearerAuth
func (s *Server) refreshHandler(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return authMiddleware.RefreshHandler
}

// logoutHandler godoc
// @Summary Logout
// @Description Logout user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/logout [post]
// @Security BearerAuth
func (s *Server) logoutHandler(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return authMiddleware.LogoutHandler
}

// healthHandler godoc
// @Summary Health check
// @Description Returns health status of the server
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
