package middleware

import (
	"errors"
	"net/http"
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

// UserIdentify to identify user token
func UserIdentify(next http.Handler) http.Handler {
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

			return []byte(helper.SigningKey), nil
		})
		if err != nil {
			helper.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
