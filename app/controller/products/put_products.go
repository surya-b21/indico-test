package products

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

// @Summary      Update Products
// @Description  Update products
// @Tags         Inventory Management
// @Accept       application/json
// @Produce		 application/json
// @Param        data   body  model.ProductAPI  true  "Body payload"
// @Router       /products/{id} [put]
// @Security BearerAuth
func PutProduct(w http.ResponseWriter, r *http.Request) {
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

	id := r.PathValue("id")
	if id == "" {
		helper.NewErrorResponse(w, http.StatusBadRequest, "please fill book id")
		return
	}

	if err := uuid.Validate(id); err == nil {
		helper.NewErrorResponse(w, http.StatusBadRequest, "id is not valid")
		return
	}

	db := service.DB
	product, err := getProduct(id, db)
	if err != nil {
		helper.NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	product.ProductAPI = body
	db.Updates(&product)

	response, _ := json.Marshal(map[string]interface{}{
		"status":  "success",
		"message": "Successfully update product",
		"data":    product,
	})

	helper.NewSuccessResponse(w, response)
}
