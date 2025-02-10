package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/internal/db"
	"backend/internal/routes"
)

func main() {
	// config.LoadConfig()
	db.InitDB()
	db.InitRedis()
	defer db.CloseRedis()
	defer db.CloseDB()

	router := gin.Default()
	router.Use(cors.Default())
	routes.RegisterUserRoutes(router)
	routes.RegisterPunchCardRoutes(router)
	router.Run(":8080")
}
