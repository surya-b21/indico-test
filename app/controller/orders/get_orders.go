package orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

// @Summary      All Orders
// @Description  Get all orders
// @Tags         Order Processing
// @Accept       application/json
// @Produce		 application/json
// @Router       /orders [get]
// @Security BearerAuth
func GetOrders(c *gin.Context) {
	db := service.DB

	orders := []model.Order{}
	db.Preload("OrderItems").Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   orders,
	})
}

// @Summary      Get Order
// @Description  Get order by id
// @Tags         Order Processing
// @Accept       application/json
// @Produce		 application/json
// @Router       /orders/{id} [get]
// @Security BearerAuth
func GetByIdOrders(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "please fill order id",
		})
		return
	}

	if err := uuid.Validate(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "id is not valid",
		})
		return
	}

	db := service.DB

	order := model.Order{}
	query := db.Preload("OrderItems").First(&order, "id = ?", id)
	if query.RowsAffected < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "order not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   order,
	})
}
