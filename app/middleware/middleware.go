package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/suryab-21/indico-test/app/helper"
)

type Middleware func(http.Handler) http.Handler

// MiddlewareStack chains multiple middlewares
func MiddlewareStack(ms ...Middleware) Middleware {
	return Middleware(func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			m := ms[i]
			next = m(next)
		}
		return next
	})
}

func ClaimToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			helper.NewErrorResponse(w, http.StatusUnauthorized, "empty auth header")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			helper.NewErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		_, err := jwt.Parse(headerParts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(os.Getenv("KEY")), nil
		})
		if err != nil {
			helper.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AdminIdentify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			helper.NewErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(headerParts[1], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(os.Getenv("KEY")), nil
		})

		if err != nil {
			helper.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if claims["role"] != "admin" {
			helper.NewErrorResponse(w, http.StatusForbidden, "Access Forbidden")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func StaffIdentify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			helper.NewErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(headerParts[1], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(os.Getenv("KEY")), nil
		})

		if err != nil {
			helper.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if claims["role"] != "staff" {
			helper.NewErrorResponse(w, http.StatusForbidden, "Access Forbidden")
			return
		}

		next.ServeHTTP(w, r)
	})
}
