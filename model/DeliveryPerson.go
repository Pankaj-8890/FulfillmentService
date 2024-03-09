package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type DeliveryPerson struct {
	gorm.Model
	ID          int64   `json:"id" gorm:"autoIncrement;primaryKey"`
	Name        string  `json:"name" gorm:"not null"`
	Availablity bool    `json:"availablity" gorm:"not null"`
	Location    string  `json:"location" gorm:"not null"`
	Orders      []Order `gorm:"foreignKey:DeliveryPersonID"`
}

type DeliveryPersonResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Availablity bool   `json:"availablity"`
	Message     string `json:"message"`
}

func (s *ServiceDb) CreateDeliveryPerson(person DeliveryPerson) (DeliveryPersonResponse, error) {

	person.Availablity = true
	res := s.DB.Create(&person)

	var deliveryPerson DeliveryPerson

	res1 := s.DB.First(&deliveryPerson)
	
	if res1.Error != nil {
        if errors.Is(res.Error, gorm.ErrRecordNotFound) {
            return DeliveryPersonResponse{}, fmt.Errorf("no delivery person found for the given conditions")
        }
        return DeliveryPersonResponse{}, res.Error
    }

	if res.RowsAffected == 0 {
		return DeliveryPersonResponse{}, errors.New("deliveryPerson creation unsuccessful")
	}

	response := DeliveryPersonResponse{
		ID:       person.ID,
		Name:     person.Name,
		Availablity: person.Availablity,
		Location: person.Location,
		Message:  "DeliveryPerson created successfuly",
	}

	return response, nil
}
