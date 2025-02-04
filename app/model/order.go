package model

import "github.com/google/uuid"

type Order struct {
	Base
	OrderAPI
	WarehouseLocation WarehouseLocation `json:"warehouse_location,omitempty"`
	OrderItems        []OrderItems      `json:"order_items,omitempty"`
}

type OrderAPI struct {
	WarehouseLocationID *uuid.UUID `json:"warehouse_location_id,omitempty" gorm:"type:varchar(36)"`
	Type                *string    `json:"type,omitempty" gorm:"type:varchar(10)"`
}
