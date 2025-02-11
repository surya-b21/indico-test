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
	"gorm.io/gorm"
)

func TestDeleteProduct(t *testing.T) {
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

	req, _ := http.NewRequest("DELETE", "/products/"+product.ID.String(), nil)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.DELETE("/products/:id", DeleteProduct)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Successfully delete product", response["message"])

	var checkProduct model.Product
	result := db.First(&checkProduct, product.ID)
	assert.Error(t, result.Error, "Product should be deleted")
	assert.True(t, gorm.ErrRecordNotFound.Error() == result.Error.Error())

	req, _ = http.NewRequest("DELETE", "/products/"+uuid.NewString(), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "product not found", response["message"])

	req, _ = http.NewRequest("DELETE", "/products/invalid-uuid", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "id is not valid", response["message"])
}
