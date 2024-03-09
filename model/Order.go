package model

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Location         string
	OrderStatus      DeliveryStatus
	DeliveryPersonID int64 `json:"DeliveryPersonId"`
}