package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/internal/config"
	"backend/internal/db"
)

func main() {
	config.LoadConfig()
	db.InitDB()
	defer db.CloseDB()

	router := gin.Default()
	router.Use(cors.Default())
	router.Run(":8080")
}
