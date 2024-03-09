package model

import "gorm.io/gorm"

type IServiceDb interface {
	CreateDeliveryPerson(DeliveryPerson) (DeliveryPersonResponse, error)
}

type ServiceDb struct {
	DB *gorm.DB
}


func NewServiceDb(db *gorm.DB) *ServiceDb {
	return &ServiceDb{DB: db}
}
