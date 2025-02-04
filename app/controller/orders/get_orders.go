package orders

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

func GetOrders(w http.ResponseWriter, r *http.Request) {
	db := service.DB

	orders := []model.Order{}
	db.Joins("OrderItems").Find(&orders)

	response, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"data":   orders,
	})

	helper.NewSuccessResponse(w, response)
}

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
