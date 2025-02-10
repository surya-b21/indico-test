package locations

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

// @Summary      Locations
// @Description  Get all locations
// @Tags         Inventory Management
// @Accept       application/json
// @Produce		 application/json
// @Router       /locations [get]
// @Security BearerAuth
func GetLocations(c *gin.Context) {
	db := service.DB

	locations := []model.WarehouseLocation{}
	db.Find(&locations)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   locations,
	})
}
