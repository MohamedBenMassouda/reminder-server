package main

import (
	"reminder-server/internal/initializers"
	"reminder-server/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	in := initializers.NewInitializers()

	router.SetupRouter(r, *in)

	r.Run(":8080")
}
