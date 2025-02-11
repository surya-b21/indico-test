package orders

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

func TestGetOrders(t *testing.T) {
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

	req, _ := http.NewRequest("GET", "/orders", nil)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/orders", GetOrders)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	assert.Empty(t, response["data"].([]interface{}))

	orderType := "receive"
	order := model.Order{
		OrderAPI: model.OrderAPI{
			WarehouseLocationID: location.ID,
			Type:                &orderType,
		},
	}
	db.Create(&order)

	orderItem := model.OrderItems{
		OrderItemsAPI: model.OrderItemsAPI{
			ProductID: product.ID,
			Quantity:  &qty1,
			OrderID:   order.ID,
		},
	}
	db.Create(&orderItem)

	req, _ = http.NewRequest("GET", "/orders", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	data := response["data"].([]interface{})
	assert.Len(t, data, 1)

	returnedOrder1 := data[0].(map[string]interface{})

	assert.NotEmpty(t, returnedOrder1["id"])

	returnedOrderItem1 := returnedOrder1["order_items"].([]interface{})[0].(map[string]interface{})

	assert.Equal(t, *orderItem.Quantity, int(returnedOrderItem1["quantity"].(float64)))
}

func TestGetByIdOrders(t *testing.T) {
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
	order := model.Order{
		OrderAPI: model.OrderAPI{
			WarehouseLocationID: location.ID,
			Type:                &orderType,
		},
	}
	db.Create(&order)

	orderItem := model.OrderItems{
		OrderItemsAPI: model.OrderItemsAPI{
			ProductID: product.ID,
			Quantity:  &qty1,
			OrderID:   order.ID,
		},
	}
	db.Create(&orderItem)

	req, _ := http.NewRequest("GET", "/orders/"+order.ID.String(), nil)
	w := httptest.NewRecorder()
	r := gin.Default()
	r.GET("/orders/:id", GetByIdOrders)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "success", response["status"])
	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	returnedOrderItems := data["order_items"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, *orderItem.Quantity, int(returnedOrderItems["quantity"].(float64)))

	req, _ = http.NewRequest("GET", "/orders/invalid-uuid", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "id is not valid", response["message"])

	req, _ = http.NewRequest("GET", "/orders/"+uuid.NewString(), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "order not found", response["message"])
}
