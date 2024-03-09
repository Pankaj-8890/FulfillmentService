package middleware

import (
	"fmt"
	"fulfillmentService/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)


func DatabaseConnection() *gorm.DB {
	host := "localhost"
	port := "5432"
	dbName := "fulfillment"
	dbUser := "postgres"
	password := "postgres"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, dbUser, dbName, password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info),})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}
	fmt.Println("Database connection successful...")

	db.Exec("DROP TABLE IF EXISTS orders;")
	db.Exec("DROP TABLE IF EXISTS delivery_people;")
	db.AutoMigrate(&model.Order{})
	db.AutoMigrate(&model.DeliveryPerson{})

	return db
}
