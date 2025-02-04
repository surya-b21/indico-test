package model

import "github.com/google/uuid"

type OrderItems struct {
	Base
	OrderItemsAPI
	Product Product
	Order   Order
}

type OrderItemsAPI struct {
	OrderID   *uuid.UUID `json:"order_id,omitempty" gorm:"type:varchar(36)"`
	ProductID *uuid.UUID `json:"product_id,omitempty" gorm:"type:varchar(36)"`
	Quantity  *int       `json:"quantity,omitempty" gorm:"type:int"`
}
