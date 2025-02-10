package auth

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

func TestSignIn(t *testing.T) {
	db := service.DBtest()
	defer func() {
		testDB, _ := db.DB()
		testDB.Close()
	}()

	os.Setenv("SALT", "test_salt") // Set env variable untuk test
	os.Setenv("KEY", "test_key")

	hash := sha256.New()
	hash.Write([]byte("password"))
	hashPassword := fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))

	body := SignInBody{
		Username: "testuser",
		Password: "password",
	}
	role := "admin"

	user := model.User{}
	user.Name = &body.Username
	user.Username = &body.Username
	user.Password = &hashPassword
	user.Role = &role

	db.Create(&user)

	jsonValue, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/signin", SignIn)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
	assert.Equal(t, float64(6*3600), data["expired"])

	body = SignInBody{
		Username: "nonexistentuser",
		Password: "password",
	}
	jsonValue, _ = json.Marshal(body)
	req, _ = http.NewRequest("POST", "/signin", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	// r.POST("/signin", SignIn)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "User not found", response["message"])
}
