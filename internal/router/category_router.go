package router

import (
	"reminder-server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRouter(router *gin.Engine, categoryHandler *handlers.CategoryHandler) {
	categories := router.Group("/categories")

	categories.GET("/", categoryHandler.List)
	categories.GET("/:id", categoryHandler.Get)

	categories.POST("/", categoryHandler.Create)

	categories.PATCH("/:id", categoryHandler.Update)

	categories.DELETE("/:id", categoryHandler.Delete)
}
