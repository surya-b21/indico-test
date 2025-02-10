package auth

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

type SignUpBody struct {
	SignInBody
	Name string `json:"name" binding:"required"`
	Role string `json:"role"`
}

// @Summary      Sign Up
// @Description  Sign up new account
// @Tags         User Authentication & Role Management
// @Accept       application/json
// @Produce		 application/json
// @Param        data   body  auth.SignUpBody  true  "Sign Up Payload"
// @Router       /register [post]
func SignUp(c *gin.Context) {
	var body SignUpBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": map[string]interface{}{
			"token":   signedToken,
			"expired": 6 * 3600,
		},
	})
}
