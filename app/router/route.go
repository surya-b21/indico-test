package router

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/suryab-21/sigmatech-test/app/middleware"
)

func InitRoutes() http.Handler {
	router := http.NewServeMux()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	authorizedMiddleware := middleware.MiddlewareStack(
		middleware.UserIdentify,
	)

	router.Handle("/", authorizedMiddleware(authorizedRoute()))

	return cors.Handler(router)
}

func authorizedRoute() *http.ServeMux {
	authorizedRoute := http.NewServeMux()

	return authorizedRoute
}
