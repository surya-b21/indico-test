package orders

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

func TestPostReceiveOrder(t *testing.T) {
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

	orderType := "receive"
	qty := 5
	body := OrderBody{
		OrderAPI: model.OrderAPI{
			WarehouseLocationID: location.ID,
			Type:                &orderType,
		},
		OrderItems: []model.OrderItemsAPI{
			{ProductID: product.ID, Quantity: &qty},
		},
	}

	jsonValue, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/orders/receive", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/orders/receive", PostReceiveOrder)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "receive order will process", response["message"])

	var updatedProduct model.Product
	db.First(&updatedProduct, product.ID)
	assert.Equal(t, 55, *updatedProduct.Quantity)
}

func TestPostShipOrder(t *testing.T) {
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

	orderType := "ship"
	qty := 5
	body := OrderBody{
		OrderAPI: model.OrderAPI{
			WarehouseLocationID: location.ID,
			Type:                &orderType,
		},
		OrderItems: []model.OrderItemsAPI{
			{ProductID: product.ID, Quantity: &qty},
		},
	}

	jsonValue, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/orders/ship", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/orders/ship", PostShipOrder)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "ship order will process", response["message"])

	time.Sleep(time.Second * 1)

	var updatedProduct model.Product
	db.First(&updatedProduct, product.ID)
	assert.Equal(t, 45, *updatedProduct.Quantity)
}
