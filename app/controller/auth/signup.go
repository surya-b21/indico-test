package auth

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

type SignUpBody struct {
	SignInBody
	Name string `json:"name"`
	Role string `json:"role"`
}

// @Summary      Sign Up
// @Description  Sign up new account
// @Tags         User Authentication & Role Management
// @Accept       application/json
// @Produce		 application/json
// @Param        data   body  auth.SignUpBody  true  "Sign Up Payload"
// @Router       /register [post]
func SignUp(w http.ResponseWriter, r *http.Request) {
	var body SignUpBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helper.NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Print(body)
	if body.Name == "" {
		helper.NewErrorResponse(w, http.StatusBadRequest, "Name is required")
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

	if body.Role == "" {
		body.Role = "staff"
	}

	body.Role = strings.ToLower(body.Role)

	hash := sha256.New()
	hash.Write([]byte(body.Password))
	hashPassword := fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))

	db := service.DB

	newUser := model.User{}
	newUser.Name = &body.Name
	newUser.Username = &body.Username
	newUser.Password = &hashPassword
	newUser.Role = &body.Role
	db.Create(&newUser)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": *newUser.ID,
		"role":    *newUser.Role,
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
