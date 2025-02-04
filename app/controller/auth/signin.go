package auth

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

type SignInBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary      Sign In
// @Description  Sign in to get bearer token
// @Tags         User Authentication & Role Management
// @Accept       application/json
// @Produce		 application/json
// @Param        data   body  auth.SignInBody  true  "Sign In Payload"
// @Router       /login [post]
func SignIn(w http.ResponseWriter, r *http.Request) {
	var body SignInBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helper.NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if body.Username == "" {
		helper.NewErrorResponse(w, http.StatusBadRequest, "Username is required")
		return
	}

	if body.Password == "" {
		helper.NewErrorResponse(w, http.StatusBadRequest, "Password is required")
		return
	}

	db := service.DB
	var user model.User
	mod := db.Where("username = ?", body.Username).First(&user)

	if mod.RowsAffected < 1 {
		helper.NewErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	hash := sha256.New()
	hash.Write([]byte(body.Password))
	hashPassword := fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))

	if hashPassword != *user.Password {
		helper.NewErrorResponse(w, http.StatusBadRequest, "Wrong user / password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": *user.ID,
		"role":    *user.Role,
		"exp":     time.Now().Add(time.Hour * 6).Unix(),
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("KEY")))
	if err != nil {
		helper.NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	response, _ := json.Marshal(map[string]interface{}{
		"token":   signedToken,
		"expired": 6 * 3600,
	})

	helper.NewSuccessResponse(w, response)
}
