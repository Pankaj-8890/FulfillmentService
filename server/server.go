package main

import (
	"context"
	"errors"
	"fmt"
	"fulfillmentService/middleware"
	"fulfillmentService/model"
	pb "fulfillmentService/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var add string = "0.0.0.0:9090"

type Server struct {
	DB *gorm.DB
	pb.FulfillmentServer
}

func main() {
	lis, err := net.Listen("tcp", add)

	if err != nil {
		log.Fatalf("Failed to Listen :%v\n", err)
	}

	log.Printf("listening %s\n", add)

	s := grpc.NewServer()
	dbClient := middleware.DatabaseConnection()
	serviceDb := model.NewServiceDb(dbClient)
	pb.RegisterFulfillmentServer(s, &Server{DB: serviceDb.DB})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to server %v\n", err)
	}
}

func (s *Server) AssignedOrder(ctx context.Context, req *pb.AssignedOrderRequest) (*pb.AssignedOrderResponse, error) {

	order := req.GetOrder()

	var deliveryPerson model.DeliveryPerson

	res := s.DB.First(&deliveryPerson, "location = ? AND availablity = ?", order.Location, true)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no delivery person available right now")
		}
		return nil, res.Error
	}

	od := model.Order{
		Location:    order.Location,
		OrderStatus: model.ASSIGNED,
	}

	deliveryPerson.Orders = append(deliveryPerson.Orders, od)
	deliveryPerson.Availablity = false
	s.DB.Save(&deliveryPerson)


	response := &pb.AssignedOrderResponse{
		Message: "Order successfully assigned",
	}

	return response, nil

}

func (s *Server) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.UpdateStatusResponse, error) {

	deliveryPersonId := req.DeliveryPersonId
	orderId := req.OrderId
	orderStatus := req.OrderStatus

	var order model.Order
	res := s.DB.Where("delivery_person_id = ? AND id = ? AND order_status NOT LIKE ?", deliveryPersonId,orderId,string(model.DELIVERED)).Find(&order)
	order.OrderStatus = model.DeliveryStatus(orderStatus)

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("no order found with this deliveryPerson")
	}

	s.DB.Save(&order)

	

	var deliveryPerson model.DeliveryPerson
	res1 := s.DB.Model(&deliveryPerson).Where("id = ?",deliveryPersonId).Update("availablity",true)

	if res1.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no deliveryPerson found")
		}
		return nil, res.Error
	}


	response := &pb.UpdateStatusResponse{
		Message:"Update order status successfully",
		OrderStatus: orderStatus,
	}

	return response,nil

}
