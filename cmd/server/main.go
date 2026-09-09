package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ticket-system/internal/config"
	"ticket-system/internal/database"
)

func main() {
	cfg := config.Load()

	database.Connect(cfg)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.Run(":" + cfg.Port)
}