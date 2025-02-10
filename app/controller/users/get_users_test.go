package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

func TestGetUsers(t *testing.T) {
	db := service.DBtest()
	defer func() {
		testDB, _ := db.DB()
		testDB.Close()
	}()

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/users", GetUsers)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	assert.Empty(t, response["data"].([]interface{}))

	user1Name := "user1"
	user2Name := "user2"
	password := "password"
	staffRole := "staff"
	adminRole := "admin"

	user1 := model.User{UserAPI: model.UserAPI{
		Name:     &user1Name,
		Username: &user1Name,
		Password: &password,
		Role:     &adminRole,
	}}
	user2 := model.User{UserAPI: model.UserAPI{
		Name:     &user2Name,
		Username: &user2Name,
		Password: &password,
		Role:     &staffRole,
	}}
	db.Create(&user1)
	db.Create(&user2)

	req, _ = http.NewRequest("GET", "/users", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	data := response["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestGetUserMe(t *testing.T) {
	db := service.DBtest()
	defer func() {
		testDB, _ := db.DB()
		testDB.Close()
	}()

	user1Name := "user1"
	password := "password"
	adminRole := "admin"

	user := model.User{UserAPI: model.UserAPI{
		Name:     &user1Name,
		Username: &user1Name,
		Password: &password,
		Role:     &adminRole,
	}}

	db.Create(&user)

	req, _ := http.NewRequest("GET", "/users/me", nil)
	w := httptest.NewRecorder()
	r := gin.Default()

	// Set user_id di context
	r.GET("/users/me", func(c *gin.Context) {
		c.Set("user_id", user.ID.String()) // Set ID user yang sesuai
		GetUserMe(c)
	})

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, *user.Name, data["name"])
	assert.Equal(t, *user.Username, data["username"])
}
