package model

import "github.com/google/uuid"

type Product struct {
	Base
	ProductAPI
	WarehouseLocation WarehouseLocation `gorm:"foreignKey:LocationID"`
}

type ProductAPI struct {
	Name       *string    `json:"name,omitempty" gorm:"type:varchar(100)"`
	Sku        *string    `json:"sku,omitempty" gorm:"type:varchar(100); unique"`
	Quantity   *int       `json:"quantity,omitempty" gorm:"type:int"`
	LocationID *uuid.UUID `json:"location_id,omitempty" gorm:"type:varchar(36)"`
}
