package orders

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
func PostReceiveOrder(c *gin.Context) {
	var body OrderBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if *body.Type == "ship" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "wrong type order",
		})
		return
	}

	db := service.DB

	if db.First(&model.WarehouseLocation{}, "id = ?", *body.WarehouseLocationID).RowsAffected < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "warehouse location not found",
		})
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

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "receive order will process",
	})
}

// @Summary      Ship Order
// @Description  Do ship order
// @Tags         Order Processing
// @Accept       application/json
// @Produce		 application/json
// @Param        data   body  orders.OrderBody  true  "Body payload"
// @Router       /orders/ship [post]
// @Security BearerAuth
func PostShipOrder(c *gin.Context) {
	var body OrderBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if *body.Type == "receive" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "wrong type order",
		})
		return
	}

	db := service.DB

	if db.First(&model.WarehouseLocation{}, "id = ?", *body.WarehouseLocationID).RowsAffected < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "warehouse location not found",
		})
		return
	}

	// check qty product
	for _, orderItem := range body.OrderItems {
		product := model.Product{}
		db.First(&product, "id = ?", orderItem.ProductID)

		if *orderItem.Quantity > *product.Quantity {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("insufficient amount for product %s", *product.Name),
			})
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

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "ship order will process",
	})
}
