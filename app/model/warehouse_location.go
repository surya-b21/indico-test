package model

type WarehouseLocation struct {
	Base
	WarehouseLocationAPI
}

type WarehouseLocationAPI struct {
	Name     *string `json:"name,omitempty" gorm:"type:varchar(100)" binding:"required"`
	Capacity *int    `json:"capacity,omitempty" gorm:"type:int" binding:"required"`
}
