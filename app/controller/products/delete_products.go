package products

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/service"
)

// @Summary      Delete Product
// @Description  Delete a product
// @Tags         Inventory Management
// @Accept       application/json
// @Produce		 application/json
// @Router       /products/{id} [delete]
// @Security BearerAuth
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
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

	db.Delete(&product)

	response, _ := json.Marshal(map[string]interface{}{
		"status":  "success",
		"message": "Successfully delete product",
	})

	helper.NewSuccessResponse(w, response)
}
