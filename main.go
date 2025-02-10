package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/suryab-21/indico-test/app/router"
	"github.com/suryab-21/indico-test/app/service"
)

// @title Indico Test
// @version 1.0
// @description For test purpose
// @host localhost:8080
// @BasePath /
// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token
func main() {
	service.InitDB()
	r := gin.Default()

	router.SetupRoutes(r)

	r.Run(":" + os.Getenv("PORT"))
}
