package main

import (
	"net/http"
	"ton-balance-api/handlers"
	"ton-balance-api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	authGroup := r.Group("/api")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.POST("/balance", handlers.GetBalance)
	}

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
