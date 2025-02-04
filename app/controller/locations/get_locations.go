package locations

import (
	"encoding/json"
	"net/http"

	"github.com/suryab-21/indico-test/app/helper"
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
func GetLocations(w http.ResponseWriter, r *http.Request) {
	db := service.DB

	locations := []model.WarehouseLocation{}
	db.Find(&locations)

	response, _ := json.Marshal(map[string]interface{}{
		"status": "success",
		"data":   locations,
	})

	helper.NewSuccessResponse(w, response)
}
