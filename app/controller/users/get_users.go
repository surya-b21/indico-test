package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

// @Summary      Users
// @Description  Get all users
// @Tags         User Authentication & Role Management
// @Accept       application/json
// @Produce		 application/json
// @Router       /users [get]
// @Security BearerAuth
func GetUsers(c *gin.Context) {
	db := service.DB

	users := []model.User{}
	db.Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   users,
	})
}

// @Summary      Users Info
// @Description  Get user info
// @Tags         User Authentication & Role Management
// @Accept       application/json
// @Produce		 application/json
// @Router       /users/me [get]
// @Security BearerAuth
func GetUserMe(c *gin.Context) {
	db := service.DB
	user := model.User{}
	db.First(&user, "id = ?", c.GetString("user_id"))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}
