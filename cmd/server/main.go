package main

import (
	"reminder-server/internal/initializers"
	"reminder-server/internal/router"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "reminder-server/docs" // swagger docs
)

// @title           Reminder Server API
// @version         1.0
// @description     A reminder management API built with Go and Gin framework
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@reminder-server.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:1323
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	r := gin.Default()

	in := initializers.NewInitializers()

	router.SetupRouter(r, *in)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + in.Config.Port)
}
