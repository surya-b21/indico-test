package locations

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

func PostLocations(w http.ResponseWriter, r *http.Request) {
	log.Println("Blocked")
	return

	var body model.WarehouseLocationAPI

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helper.NewErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if body.Name == nil {
		helper.NewErrorResponse(w, http.StatusUnprocessableEntity, "Name is required")
		return
	}

	if body.Capacity == nil {
		helper.NewErrorResponse(w, http.StatusUnprocessableEntity, "Capacity is required")
		return
	}

	db := service.DB
	location := model.WarehouseLocation{}
	location.WarehouseLocationAPI = body
	db.Create(&location)

	response, _ := json.Marshal(map[string]interface{}{
		"message": "success",
		"data":    location,
	})

	helper.NewSuccessResponse(w, response)
}
