package products

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
	"gorm.io/gorm"
)

// @Summary      All Products
// @Description  Get all products
// @Tags         Inventory Management
// @Accept       application/json
// @Produce		 application/json
// @Router       /products [get]
// @Security BearerAuth
func GetProducts(w http.ResponseWriter, r *http.Request) {
	db := service.DB

	products := []model.Product{}
	db.Joins("WarehouseLocation").Find(&products)

	response, _ := json.Marshal(map[string]interface{}{
		"message": "success",
		"data":    products,
	})

	helper.NewSuccessResponse(w, response)
}

// @Summary      Products
// @Description  Get products by id
// @Tags         Inventory Management
// @Accept       application/json
// @Produce		 application/json
// @Router       /products/{id} [get]
// @Security BearerAuth
func GetByIdProducts(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		helper.NewErrorResponse(w, http.StatusBadRequest, "please fill product id")
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

	response, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"data":   *product,
	})

	helper.NewSuccessResponse(w, response)
}

func getProduct(id string, db *gorm.DB) (*model.Product, error) {
	product := model.Product{}
	query := db.First(&product, "id = ?", id)

	if query.RowsAffected < 1 {
		return nil, errors.New("product not found")
	}

	return &product, nil
}
