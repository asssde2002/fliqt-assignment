package routes

import (
	"backend/internal/handlers"
	"backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/signup", handlers.SignUp)
	publicRoutes.POST("/login", handlers.Login)

	privateRoutes := router.Group("user")
	privateRoutes.Use(middlewares.AuthenticationMiddleware())
	privateRoutes.GET("/:id", handlers.GetUser)
	privateRoutes.PUT("/:id/roles", handlers.PutUserRoles)
}
