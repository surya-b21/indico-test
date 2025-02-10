package products

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

// @Summary      Post Products
// @Description  Add new products
// @Tags         Inventory Management
// @Accept       application/json
// @Produce		 application/json
// @Param        data   body  model.ProductAPI  true  "Body payload"
// @Router       /products [post]
// @Security BearerAuth
func PostProduct(c *gin.Context) {
	var body model.ProductAPI
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	db := service.DB
	product := model.Product{}
	product.ProductAPI = body
	db.Create(&product)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully add new product",
		"data":    product,
	})
}
