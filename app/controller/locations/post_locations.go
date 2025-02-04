package locations

import (
	"encoding/json"
	"net/http"

	"github.com/suryab-21/indico-test/app/helper"
	"github.com/suryab-21/indico-test/app/model"
	"github.com/suryab-21/indico-test/app/service"
)

func PostLocations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.NewErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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
		"status":  "success",
		"message": "Successfully add new location",
		"data":    location,
	})

	helper.NewSuccessResponse(w, response)
}
