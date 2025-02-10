package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/suryab-21/indico-test/app/service"
)

func TestSignUp(t *testing.T) {
	db := service.DBtest()
	defer func() {
		testDB, _ := db.DB()
		testDB.Close()
	}()

	os.Setenv("SALT", "test_salt")
	os.Setenv("KEY", "test_key")

	body := SignUpBody{}
	body.Username = "staffuser"
	body.Password = "password"
	body.Name = "Staff User"

	jsonValue, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/register", SignUp)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "success", response["status"])
	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
	assert.Equal(t, float64(6*3600), data["expired"])

	body = SignUpBody{
		SignInBody: SignInBody{
			Username: "anotheruser",
			Password: "password",
		},
		Name: "Another User",
		Role: "admin",
	}

	jsonValue, _ = json.Marshal(body)
	req, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	data = response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
	assert.Equal(t, float64(6*3600), data["expired"])
}
