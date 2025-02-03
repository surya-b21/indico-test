package model

type WarehouseLocation struct {
	Base
}

type WarehouseLocationAPI struct {
	Name     *string `json:"name,omitempty" gorm:"type:varchar(100)"`
	Capacity *int    `json:"capacity,omitempty" gorm:"type:int"`
}
