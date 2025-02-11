package locations

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

func TestGetLocations(t *testing.T) {
	db := service.DBtest()
	defer func() {
		testDB, _ := db.DB()
		testDB.Close()
	}()

	req, _ := http.NewRequest("GET", "/locations", nil)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/locations", GetLocations)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	assert.Empty(t, response["data"].([]interface{}))

	location1Name := "Location 1"
	location2Name := "Location 2"
	capacity := 100

	location1 := model.WarehouseLocation{
		WarehouseLocationAPI: model.WarehouseLocationAPI{
			Name:     &location1Name,
			Capacity: &capacity,
		},
	}

	location2 := model.WarehouseLocation{
		WarehouseLocationAPI: model.WarehouseLocationAPI{
			Name:     &location2Name,
			Capacity: &capacity,
		},
	}

	db.Create(&location1)
	db.Create(&location2)

	req, _ = http.NewRequest("GET", "/locations", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	data := response["data"].([]interface{})
	assert.Len(t, data, 2)

	returnedLocation1 := data[0].(map[string]interface{})
	returnedLocation2 := data[1].(map[string]interface{})

	assert.Equal(t, *location1.Name, returnedLocation1["name"])
	assert.Equal(t, *location2.Name, returnedLocation2["name"])

}
