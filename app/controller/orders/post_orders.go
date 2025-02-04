package orders

import (
	"encoding/json"
	"net/http"

	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

type OrderBody struct {
	model.OrderAPI
	OrderItems []model.OrderItemsAPI `json:"order_items,omitempty"`
}

func PostReceiveOrder(w http.ResponseWriter, r *http.Request) {
	var body OrderBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helper.NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if *body.Type == "ship" {
		helper.NewErrorResponse(w, http.StatusBadRequest, "wrong type order")
		return
	}

	db := service.DB

	if db.First(&model.WarehouseLocation{}, "id = ?", *body.WarehouseLocationID).RowsAffected < 1 {
		helper.NewErrorResponse(w, http.StatusNotFound, "warehouse location not found")
		return
	}

	go func() {
		order := model.Order{}
		order.WarehouseLocationID = body.WarehouseLocationID
		order.Type = body.Type
		db.Create(&order)

		orderItems := []model.OrderItems{}
		for _, orderItem := range body.OrderItems {
			orderItems = append(orderItems, model.OrderItems{
				OrderItemsAPI: orderItem,
			})
		}

		db.CreateInBatches(&orderItems, 100)
	}()

	response, _ := json.Marshal(map[string]interface{}{
		"status":  "success",
		"message": "receive order will process",
	})

	helper.NewSuccessResponse(w, response)
}

func PostShipOrder(w http.ResponseWriter, r *http.Request) {
	var body OrderBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helper.NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if *body.Type == "receive" {
		helper.NewErrorResponse(w, http.StatusBadRequest, "wrong type order")
		return
	}

	db := service.DB

	if db.First(&model.WarehouseLocation{}, "id = ?", *body.WarehouseLocationID).RowsAffected < 1 {
		helper.NewErrorResponse(w, http.StatusNotFound, "warehouse location not found")
		return
	}

	go func() {
		order := model.Order{}
		order.WarehouseLocationID = body.WarehouseLocationID
		order.Type = body.Type
		db.Create(&order)

		orderItems := []model.OrderItems{}
		for _, orderItem := range body.OrderItems {
			orderItems = append(orderItems, model.OrderItems{
				OrderItemsAPI: orderItem,
			})
		}

		db.CreateInBatches(&orderItems, 100)
	}()

	response, _ := json.Marshal(map[string]interface{}{
		"status":  "success",
		"message": "receive order will process",
	})

	helper.NewSuccessResponse(w, response)
}
