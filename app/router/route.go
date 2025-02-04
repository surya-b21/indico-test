package router

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/suryab-21/indico-test/app/controller/auth"
	"github.com/suryab-21/indico-test/app/controller/locations"
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

	authorizedRoute.HandleFunc("GET /users/me", users.GetUserMe)
	authorizedRoute.HandleFunc("GET /locations", locations.GetLocations)

	// admin route
	authorizedRoute.Handle("/users", middleware.AdminIdentify(http.HandlerFunc(users.GetUsers)))
	authorizedRoute.Handle("/locations", middleware.AdminIdentify(http.HandlerFunc(locations.PostLocations)))

	return authorizedRoute
}
