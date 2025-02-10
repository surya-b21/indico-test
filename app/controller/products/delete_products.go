package products

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suryab-21/indico-test/app/service"
)

// @Summary      Delete Product
// @Description  Delete a product
// @Tags         Inventory Management
// @Accept       application/json
// @Produce		 application/json
// @Router       /products/{id} [delete]
// @Security BearerAuth
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "please fill book id",
		})
		return
	}

	if err := uuid.Validate(id); err == nil {
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

	db.Delete(&product)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully delete product",
	})
}
