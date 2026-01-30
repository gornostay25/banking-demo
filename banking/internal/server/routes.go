package server

import (
	"log"
	"net/http"

	"banking/internal/services"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
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
		if ctx.Request.URL.Path == "/swagger" {
			ctx.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
			return
		}
		ctx.Next()
	}, ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

// loginHandler godoc
// @Summary Login
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body services.LoginRequest true "Login credentials"
// @Success 200 {object} services.LoginResponse
// @Failure 401 {object} services.ErrorResponse
// @Router /api/auth/login [post]
func (s *Server) loginHandler(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return authMiddleware.LoginHandler
}

// refreshHandler godoc
// @Summary Refresh token
// @Description Refresh JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} services.LoginResponse
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
