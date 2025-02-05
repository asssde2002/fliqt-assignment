package routes

import (
	"backend/internal/handlers"
	"backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterPunchCardRoutes(router *gin.Engine) {
	privateRoutes := router.Group("/clock")
	privateRoutes.Use(middlewares.AuthenticationMiddleware())
	privateRoutes.POST("/in", handlers.ClockIn)
	privateRoutes.POST("/out", handlers.ClockOut)
}
