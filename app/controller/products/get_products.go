package products

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
func GetProducts(c *gin.Context) {
	db := service.DB

	products := []model.Product{}
	db.Joins("WarehouseLocation").Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    products,
	})
}

// @Summary      Products
// @Description  Get products by id
// @Tags         Inventory Management
// @Accept       application/json
// @Produce		 application/json
// @Router       /products/{id} [get]
// @Security BearerAuth
func GetByIdProducts(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "please fill product id",
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
	product, err := getProduct(id, db)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   *product,
	})
}

func getProduct(id string, db *gorm.DB) (*model.Product, error) {
	product := model.Product{}
	query := db.First(&product, "id = ?", id)

	if query.RowsAffected < 1 {
		return nil, errors.New("product not found")
	}

	return &product, nil
}
