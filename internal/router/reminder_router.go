package router

import (
	"reminder-server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupReminderRouter(router *gin.Engine, reminderHandler *handlers.ReminderHandler) {
	reminders := router.Group("/reminders")

	reminders.GET("/", reminderHandler.List)
	reminders.GET("/:id", reminderHandler.Get)

	reminders.POST("/", reminderHandler.Create)

	reminders.PATCH("/:id", reminderHandler.Update)

	reminders.DELETE("/:id", reminderHandler.Delete)
}
