package locations

import (
	"encoding/json"
	"net/http"

	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

func GetLocations(w http.ResponseWriter, r *http.Request) {
	db := service.DB

	locations := []model.WarehouseLocation{}
	db.Find(&locations)

	response, _ := json.Marshal(map[string]interface{}{
		"message": "success",
		"data":    locations,
	})

	helper.NewSuccessResponse(w, response)
}
