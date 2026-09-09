package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ticket-system/internal/config"
	"ticket-system/internal/database"
	"ticket-system/internal/handlers"
	"ticket-system/internal/middleware"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	authHandler := handlers.NewAuthHandler(db, cfg)

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)

	protected := router.Group("/tickets")
	protected.Use(middleware.JWTAuth(cfg))

	// Ticket routes will be added here.

	router.Run(":" + cfg.Port)
}