package locations

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

func TestPostLocations(t *testing.T) {
	db := service.DBtest()
	defer func() {
		testDB, _ := db.DB()
		testDB.Close()
	}()

	location1Name := "Location 1"
	capacity := 100

	body := model.WarehouseLocationAPI{
		Name:     &location1Name,
		Capacity: &capacity,
	}

	jsonValue, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/locations", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/locations", PostLocations)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Successfully add new location", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, *body.Name, data["name"])
	assert.NotEmpty(t, data["id"])

}
