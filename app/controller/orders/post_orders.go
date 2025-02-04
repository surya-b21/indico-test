package orders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

type OrderBody struct {
	model.OrderAPI
	OrderItems []model.OrderItemsAPI `json:"order_items,omitempty"`
}

// @Summary      Receive Order
// @Description  Do receive order
// @Tags         Order Processing
// @Accept       application/json
// @Produce		 application/json
// @Param        data   body  orders.OrderBody  true  "Body payload"
// @Router       /orders/receive [post]
// @Security BearerAuth
func PostReceiveOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.NewErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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
				OrderItemsAPI: model.OrderItemsAPI{
					OrderID:   order.ID,
					ProductID: orderItem.ProductID,
					Quantity:  orderItem.Quantity,
				},
			})

			product := model.Product{}
			db.First(&product, "id = ?", orderItem.ProductID)

			newQty := *product.Quantity + *orderItem.Quantity
			product.Quantity = &newQty
			db.Save(&product)
		}

		db.CreateInBatches(&orderItems, 100)

	}()

	response, _ := json.Marshal(map[string]interface{}{
		"status":  "success",
		"message": "receive order will process",
	})

	helper.NewSuccessResponse(w, response)
}

// @Summary      Ship Order
// @Description  Do ship order
// @Tags         Order Processing
// @Accept       application/json
// @Produce		 application/json
// @Param        data   body  orders.OrderBody  true  "Body payload"
// @Router       /orders/ship [post]
// @Security BearerAuth
func PostShipOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.NewErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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

	// check qty product
	for _, orderItem := range body.OrderItems {
		product := model.Product{}
		db.First(&product, "id = ?", orderItem.ProductID)

		if *orderItem.Quantity > *product.Quantity {
			helper.NewErrorResponse(w, http.StatusNotFound, fmt.Sprintf("insufficient amount for product %s", *product.Name))
			return
		}
	}

	go func() {
		order := model.Order{}
		order.WarehouseLocationID = body.WarehouseLocationID
		order.Type = body.Type
		db.Create(&order)

		orderItems := []model.OrderItems{}
		for _, orderItem := range body.OrderItems {
			orderItems = append(orderItems, model.OrderItems{
				OrderItemsAPI: model.OrderItemsAPI{
					OrderID:   order.ID,
					ProductID: orderItem.ProductID,
					Quantity:  orderItem.Quantity,
				},
			})

			product := model.Product{}
			db.First(&product, "id = ?", orderItem.ProductID)

			newQty := *product.Quantity - *orderItem.Quantity
			product.Quantity = &newQty
			db.Save(&product)
		}

		db.CreateInBatches(&orderItems, 100)
	}()

	response, _ := json.Marshal(map[string]interface{}{
		"status":  "success",
		"message": "receive order will process",
	})

	helper.NewSuccessResponse(w, response)
}
