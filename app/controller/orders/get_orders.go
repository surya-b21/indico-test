package orders

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/suryab-21/indico-test/app/helper"
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
func GetOrders(w http.ResponseWriter, r *http.Request) {
	db := service.DB

	orders := []model.Order{}
	db.Preload("OrderItems").Find(&orders)

	response, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"data":   orders,
	})

	helper.NewSuccessResponse(w, response)
}

// @Summary      Get Order
// @Description  Get order by id
// @Tags         Order Processing
// @Accept       application/json
// @Produce		 application/json
// @Router       /orders/{id} [get]
// @Security BearerAuth
func GetByIdOrders(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		helper.NewErrorResponse(w, http.StatusBadRequest, "please fill order id")
		return
	}

	if err := uuid.Validate(id); err == nil {
		helper.NewErrorResponse(w, http.StatusBadRequest, "id is not valid")
		return
	}

	db := service.DB

	order := model.Order{}
	query := db.Joins("OrderItems").First(&order, "id = ?", id)
	if query.RowsAffected < 1 {
		helper.NewErrorResponse(w, http.StatusNotFound, "order not found")
	}

	response, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"data":   order,
	})

	helper.NewSuccessResponse(w, response)
}
