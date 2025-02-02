package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/handlers"
)

func main() {
	config.LoadConfig()
	db.InitDB()
	defer db.CloseDB()

	router := gin.Default()
	router.Use(cors.Default())
	router.POST("auth/signup", handlers.SignUp)
	router.GET("user/:id", handlers.GetUser)
	router.Run(":8080")
}
