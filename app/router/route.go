package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/suryab-21/indico-test/app/controller/auth"
	"github.com/suryab-21/indico-test/app/controller/locations"
	"github.com/suryab-21/indico-test/app/controller/orders"
	"github.com/suryab-21/indico-test/app/controller/products"
	"github.com/suryab-21/indico-test/app/controller/users"
	"github.com/suryab-21/indico-test/app/middleware"
)

func SetupRoutes(r *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Authorization", "Accept"}

	r.Use(cors.New(config))

	r.POST("/register", auth.SignUp)
	r.POST("/login", auth.SignIn)

	r.Use(middleware.ClaimToken())

	// users
	r.GET("/users/me", users.GetUserMe)
	r.GET("/users", middleware.AdminIdentify(), users.GetUsers)

	// locations
	r.GET("/locations", locations.GetLocations)
	r.POST("/locations", middleware.AdminIdentify(), locations.PostLocations)

	// products
	r.GET("/products", products.GetProducts)
	r.POST("/products", middleware.AdminIdentify(), products.PostProduct)
	r.GET("/products/:id", products.GetByIdProducts)
	r.PUT("/products/:id", middleware.AdminIdentify(), products.PutProduct)
	r.DELETE("/products/:id", middleware.AdminIdentify(), products.DeleteProduct)

	// orders
	r.GET("/orders", orders.GetOrders)
	r.POST("/orders/receive", orders.PostReceiveOrder)
	r.POST("/orders/ship", orders.PostShipOrder)
	r.GET("/orders/:id", orders.GetByIdOrders)
}
