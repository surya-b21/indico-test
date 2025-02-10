package model

import "github.com/google/uuid"

type OrderItems struct {
	Base
	OrderItemsAPI
	Product Product `json:"product,omitempty"`
	Order   Order   `json:"order,omitempty"`
}

type OrderItemsAPI struct {
	OrderID   *uuid.UUID `json:"order_id,omitempty" gorm:"type:varchar(36)" binding:"required"`
	ProductID *uuid.UUID `json:"product_id,omitempty" gorm:"type:varchar(36)" binding:"required"`
	Quantity  *int       `json:"quantity,omitempty" gorm:"type:int" binding:"required"`
}
