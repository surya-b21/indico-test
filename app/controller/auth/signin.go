package auth

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

type SignInBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary      Sign In
// @Description  Sign in to get bearer token
// @Tags         User Authentication & Role Management
// @Accept       application/json
// @Produce		 application/json
// @Param        data   body  auth.SignInBody  true  "Sign In Payload"
// @Router       /login [post]
func SignIn(c *gin.Context) {
	var body SignInBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	db := service.DB
	var user model.User
	mod := db.Where("username = ?", body.Username).First(&user)

	if mod.RowsAffected < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	hash := sha256.New()
	hash.Write([]byte(body.Password))
	hashPassword := fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))

	if hashPassword != *user.Password {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Wrong username or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": *user.ID,
		"role":    *user.Role,
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
