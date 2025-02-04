package products

import (
	"encoding/json"
	"net/http"

	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

func PostProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.NewErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var body model.ProductAPI
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helper.NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if body.Name == nil {
		helper.NewErrorResponse(w, http.StatusUnprocessableEntity, "name is required")
		return
	}

	if body.Sku == nil {
		helper.NewErrorResponse(w, http.StatusUnprocessableEntity, "sku is required")
		return
	}

	if body.Quantity == nil {
		helper.NewErrorResponse(w, http.StatusUnprocessableEntity, "quantity is required")
		return
	}

	if body.LocationID == nil {
		helper.NewErrorResponse(w, http.StatusUnprocessableEntity, "location id is required")
		return
	}

	db := service.DB
	product := model.Product{}
	product.ProductAPI = body
	db.Create(&product)

	response, _ := json.Marshal(map[string]interface{}{
		"status":  "success",
		"message": "Successfully add new product",
		"data":    product,
	})

	helper.NewSuccessResponse(w, response)
}
