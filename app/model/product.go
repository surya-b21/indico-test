package model

import "github.com/google/uuid"

type Product struct {
	Base
	ProductAPI
	WarehouseLocation WarehouseLocation `json:"warehouse_location,omitempty" gorm:"foreignKey:LocationID"`
}

type ProductAPI struct {
	Name       *string    `json:"name,omitempty" gorm:"type:varchar(100)" binding:"required"`
	Sku        *string    `json:"sku,omitempty" gorm:"type:varchar(100); unique" binding:"required"`
	Quantity   *int       `json:"quantity,omitempty" gorm:"type:int" binding:"required"`
	LocationID *uuid.UUID `json:"location_id,omitempty" gorm:"type:varchar(36)" binding:"required"`
}
