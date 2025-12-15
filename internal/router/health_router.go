package router

import (
	"reminder-server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupHealthRouter(router *gin.Engine, healthHandler *handlers.HealthHandler) {
	health := router.Group("/health")
	{
		health.GET("/liveness", healthHandler.Liveness)
		health.GET("/readiness", healthHandler.Readiness)
		health.GET("/startup", healthHandler.Startup)
	}
}
