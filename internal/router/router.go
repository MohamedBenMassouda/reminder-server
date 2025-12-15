package router

import (
	"reminder-server/internal/initializers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, in initializers.Initializers) *gin.Engine {
	SetupHealthRouter(router, in.HealthHandler)
	SetupCategoryRouter(router, in.CategoryHandler)
	SetupReminderRouter(router, in.ReminderHandler)
	SetupUserRouter(router, in.UserHandler)

	return router
}
