package controller

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	model "fulfillmentService/model"
)



func InitController(r *mux.Router, serviceDb model.IServiceDb) {

	r.HandleFunc("/deliveryPersons", createDeliveryPerson(serviceDb)).Methods("POST")
	
}

func createDeliveryPerson(serviceDb model.IServiceDb) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {


		w.Header().Set("Content-Type", "application/json")
		var person model.DeliveryPerson
		json.NewDecoder(r.Body).Decode(&person)

		deliveryPerson,err := serviceDb.CreateDeliveryPerson(person)

		if err!=nil{
			json.NewEncoder(w).Encode(err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(deliveryPerson)

	}
}





