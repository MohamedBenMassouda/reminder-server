package router

import (
	"reminder-server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(router *gin.Engine, userHandler *handlers.UserHandler) {
	users := router.Group("/users")

	users.GET("/", userHandler.List)
	users.GET("/:id", userHandler.Get)

	users.POST("/signup", userHandler.SignUp)
	users.POST("/login", userHandler.Login)
}
