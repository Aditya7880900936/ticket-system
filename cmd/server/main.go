package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ticket-system/internal/config"
	"ticket-system/internal/database"
	"ticket-system/internal/handlers"
	"ticket-system/internal/middleware"
	"ticket-system/internal/repository"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://ticket-system-frontend-9rf1.onrender.com",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PATCH",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: false,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	authHandler := handlers.NewAuthHandler(db, cfg)

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)

	ticketRepository := repository.NewTicketRepository(db)
	ticketHandler := handlers.NewTicketHandler(ticketRepository)

	protected := router.Group("/tickets")
	protected.Use(middleware.JWTAuth(cfg))

	protected.POST("", ticketHandler.Create)
	protected.GET("", ticketHandler.List)
	protected.GET("/:id", ticketHandler.GetByID)
	protected.PATCH("/:id/status", ticketHandler.UpdateStatus)

	router.Run(":" + cfg.Port)
}
