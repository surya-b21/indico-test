package products

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

func TestPostProduct(t *testing.T) {
	db := service.DBtest()
	defer func() {
		testDB, _ := db.DB()
		testDB.Close()
	}()

	locationName := "Location Test"
	location := model.WarehouseLocation{
		WarehouseLocationAPI: model.WarehouseLocationAPI{
			Name: &locationName,
		},
	}

	db.Create(&location)

	product1Name := "Product 1"
	sku1 := "000001"
	qty1 := 50
	body := model.ProductAPI{
		Name:       &product1Name,
		Sku:        &sku1,
		Quantity:   &qty1,
		LocationID: location.ID,
	}

	jsonValue, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/products", PostProduct)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Successfully add new product", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, *body.Name, data["name"])
	assert.Equal(t, body.LocationID.String(), data["location_id"])
	assert.NotEmpty(t, data["id"])
}
