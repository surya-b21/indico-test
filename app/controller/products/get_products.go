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

func GetProducts(w http.ResponseWriter, r *http.Request) {
	db := service.DB

	products := []model.Product{}
	db.Find(&products)

	response, _ := json.Marshal(map[string]interface{}{
		"message": "success",
		"data":    products,
	})

	helper.NewSuccessResponse(w, response)
}

func GetByIdProducts(w http.ResponseWriter, r *http.Request) {
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
