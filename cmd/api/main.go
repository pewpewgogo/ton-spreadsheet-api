package main

import (
	"ton-balance-api/handlers"
	"ton-balance-api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

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
