package server

import (
	"context"
	"errors"
	"fmt"
	pb "fulfillmentService/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var add string = "0.0.0.0:9090"

type Server struct {
	DB *gorm.DB
	pb.FulfillmentServer
}

type DeliveryPerson struct {
	gorm.Model
	ID       int64  `json:"id" gorm:"autoIncrement"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func DatabaseConnection() {
	host := "localhost"
	port := "5432"
	dbName := "fulfillment"
	dbUser := "postgres"
	password := "postgres"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, dbUser, dbName, password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}
	fmt.Println("Database connection successful...")

	// Auto-migrate models
	db.AutoMigrate(&DeliveryPerson{})

}

func Init() {
	lis, err := net.Listen("tcp", add)

	if err != nil {
		log.Fatalf("Failed to Listen :%v\n", err)
	}

	log.Printf("listening %s\n", add)

	s := grpc.NewServer()
	pb.RegisterFulfillmentServer(s, &Server{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to server %v\n", err)
	}
}

func (s *Server) CreateDeliveryPerson(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {

	deliveryPerson := DeliveryPerson{
		Name: req.DeliveryPerson.Name,
		Location: req.DeliveryPerson.Location,
	}

	res := s.DB.Create(&deliveryPerson)

	if res.RowsAffected == 0 {
		return nil, errors.New("delivery person creation failed")
	}

	d := &pb.DeliveryPerson{
		Name: deliveryPerson.Name,
		Location: deliveryPerson.Location,
	}

	response := &pb.CreateResponse{
		DeliveryPerson: d,
		Message:        "Delivery person successfully created",
	}

	return response, nil
}
