package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.NewErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	db := service.DB

	users := []model.User{}
	db.Find(&users)

	response, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"data":   users,
	})

	helper.NewSuccessResponse(w, response)
}

func GetUserMe(w http.ResponseWriter, r *http.Request) {
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

	db := service.DB
	user := model.User{}
	db.First(&user, "id = ?", claims["user_id"])

	response, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"data":   user,
	})

	helper.NewSuccessResponse(w, response)
}
