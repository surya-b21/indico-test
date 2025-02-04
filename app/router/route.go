package router

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/suryab-21/indico-test/app/controller/auth"
	"github.com/suryab-21/indico-test/app/controller/locations"
	"github.com/suryab-21/indico-test/app/controller/orders"
	"github.com/suryab-21/indico-test/app/controller/products"
	"github.com/suryab-21/indico-test/app/controller/users"
	"github.com/suryab-21/indico-test/app/middleware"
)

func InitRoutes() http.Handler {
	router := http.NewServeMux()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	authorizedMiddleware := middleware.MiddlewareStack(
		middleware.ClaimToken,
	)

	router.Handle("/", authorizedMiddleware(authorizedRoute()))
	router.HandleFunc("POST /register", auth.SignUp)
	router.HandleFunc("POST /login", auth.SignIn)

	return cors.Handler(router)
}

func authorizedRoute() *http.ServeMux {
	authorizedRoute := http.NewServeMux()

	// users
	authorizedRoute.HandleFunc("GET /users/me", users.GetUserMe)
	authorizedRoute.Handle("/users", middleware.AdminIdentify(http.HandlerFunc(users.GetUsers)))

	// locations
	authorizedRoute.HandleFunc("GET /locations", locations.GetLocations)
	authorizedRoute.Handle("/locations", middleware.AdminIdentify(http.HandlerFunc(locations.PostLocations)))

	// products
	authorizedRoute.HandleFunc("GET /products", products.GetProducts)
	authorizedRoute.Handle("/products", middleware.AdminIdentify(http.HandlerFunc(products.PostProduct)))
	authorizedRoute.HandleFunc("GET /products/:id", products.GetByIdProducts)
	authorizedRoute.Handle("/products/:id", middleware.AdminIdentify(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			products.PutProduct(w, r)
		case http.MethodDelete:
			products.DeleteProduct(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	// orders
	authorizedRoute.HandleFunc("GET /orders", orders.GetOrders)
	authorizedRoute.Handle("/orders/receive", middleware.AdminIdentify(http.HandlerFunc(orders.PostReceiveOrder)))
	authorizedRoute.Handle("/orders/ship", middleware.AdminIdentify(http.HandlerFunc(orders.PostShipOrder)))
	authorizedRoute.HandleFunc("GET /orders/:id", orders.GetByIdOrders)

	return authorizedRoute
}
