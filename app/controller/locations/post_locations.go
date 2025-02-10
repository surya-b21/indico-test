package locations

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

// @Summary      Add Locations
// @Description  Add new locations
// @Tags         Inventory Management
// @Accept       application/json
// @Produce		 application/json
// @Param        data   body  model.WarehouseLocationAPI  true  "Body payload"
// @Router       /locations [post]
// @Security BearerAuth
func PostLocations(c *gin.Context) {
	var body model.WarehouseLocationAPI

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	db := service.DB
	location := model.WarehouseLocation{}
	location.WarehouseLocationAPI = body
	db.Create(&location)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully add new location",
		"data":    location,
	})
}
