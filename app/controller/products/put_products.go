package products

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
func PutProduct(c *gin.Context) {
	var body model.ProductAPI
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "please fill book id",
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

	product.ProductAPI = body
	db.Updates(&product)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully update product",
		"data":    product,
	})
}
