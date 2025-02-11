package products

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

func TestGetProducts(t *testing.T) {
	db := service.DBtest()
	defer func() {
		testDB, _ := db.DB()
		testDB.Close()
	}()

	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/products", GetProducts)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["message"])
	assert.Empty(t, response["data"].([]interface{}))

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
	product1 := model.Product{
		ProductAPI: model.ProductAPI{
			Name:       &product1Name,
			Sku:        &sku1,
			Quantity:   &qty1,
			LocationID: location.ID,
		},
	}
	db.Create(&product1)

	req, _ = http.NewRequest("GET", "/products", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["message"])
	data := response["data"].([]interface{})
	assert.Len(t, data, 1)

	returnedProduct1 := data[0].(map[string]interface{})

	assert.Equal(t, *product1.Name, returnedProduct1["name"])
	assert.Equal(t, *location.Name, returnedProduct1["warehouse_location"].(map[string]interface{})["name"])
}

func TestGetByIdProducts(t *testing.T) {
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
	product := model.Product{
		ProductAPI: model.ProductAPI{
			Name:       &product1Name,
			Sku:        &sku1,
			Quantity:   &qty1,
			LocationID: location.ID,
		},
	}
	db.Create(&product)

	req, _ := http.NewRequest("GET", "/products/"+product.ID.String(), nil)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/products/:id", GetByIdProducts)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, *product.Name, data["name"])

	req, _ = http.NewRequest("GET", "/products/invalid-uuid", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "id is not valid", response["message"])

	req, _ = http.NewRequest("GET", "/products/"+uuid.NewString(), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "product not found", response["message"])
}
