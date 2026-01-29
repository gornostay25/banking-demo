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

	// r.GET("/websocket", s.websocketHandler)

	// AUTH
	authService := services.NewAuthService(s.db)
	authMiddleware, err := jwt.New(authService.JWTInitParams())
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	authGroup := public.Group("/auth")
	{
		authGroup.POST("/login", authService.LoginHandler)
		authGroup.GET("/me", authMiddleware.MiddlewareFunc(), authService.MeHandler)
	}

	protectedGroup := public.Group("")
	protectedGroup.Use(authMiddleware.MiddlewareFunc())

	// Account Operations

	accountGroup := protectedGroup.Group("/accounts")
	{
		accountGroup.GET("", s.HelloWorldHandler)
		accountGroup.GET("/:id/balance", s.GetAccountBalanceHandler)
	}

	// Transaction Operations

	transactionGroup := protectedGroup.Group("/transactions")
	{
		transactionGroup.POST("/transfer", s.TransferHandler)
		transactionGroup.POST("/exchange", s.ExchangeHandler)
		transactionGroup.GET("", s.GetTransactionsHandler)
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

// HelloWorldHandler godoc
// @Summary Hello World handler
// @Description Returns a hello world message
// @Tags accounts
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/accounts [get]
// @Security BearerAuth
func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
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

// GetAccountBalanceHandler godoc
// @Summary Get account balance
// @Description Returns balance for a specific account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} map[string]string
// @Router /api/accounts/{id}/balance [get]
// @Security BearerAuth
func (s *Server) GetAccountBalanceHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

// TransferHandler godoc
// @Summary Transfer money
// @Description Transfer money between accounts
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/transactions/transfer [post]
// @Security BearerAuth
func (s *Server) TransferHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

// ExchangeHandler godoc
// @Summary Exchange currency
// @Description Exchange currency between accounts
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/transactions/exchange [post]
// @Security BearerAuth
func (s *Server) ExchangeHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

// GetTransactionsHandler godoc
// @Summary Get transaction history
// @Description Returns transaction history
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/transactions [get]
// @Security BearerAuth
func (s *Server) GetTransactionsHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

// func (s *Server) websocketHandler(c *gin.Context) {
// 	w := c.Writer
// 	r := c.Request
// 	socket, err := websocket.Accept(w, r, nil)
// 	if err != nil {
// 		log.Printf("could not open websocket: %v", err)
// 		_, _ = w.Write([]byte("could not open websocket"))
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	defer socket.Close(websocket.StatusGoingAway, "server closing websocket")

// 	ctx := r.Context()
// 	socketCtx := socket.CloseRead(ctx)

// 	for {
// 		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
// 		err := socket.Write(socketCtx, websocket.MessageText, []byte(payload))
// 		if err != nil {
// 			break
// 		}
// 		time.Sleep(time.Second * 2)
// 	}
// }
