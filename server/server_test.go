package main

import (
	"context"
	"encoding/base64"
	"testing"

	pb "fulfillmentService/proto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateDeliveryExecutive(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error creating mock db: %v", err)
	defer db.Close()

	dialect := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialect, &gorm.Config{})
	assert.Nil(t, err, "Error creating mock gorm db: %v", err)

	tests := []struct {
		name     string
		req      *pb.CreateDeliveryExecutiveReuqest
		rows     func()
		want     *pb.CreateDeliveryExecutiveResponse
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Create delivery executive - Success",
			req: &pb.CreateDeliveryExecutiveReuqest{
				DeliveryExecutive: &pb.Deliver_Executive{
					Name:     "John Doe",
					Location: "New York",
					Password: "password123",
				},
			},
			rows: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "name", "location", "password"}).AddRow(1, "John Doe", "New York", "password123")
				mock.ExpectQuery("INSERT").WillReturnRows(rows)
				mock.ExpectCommit()
			},
			want: &pb.CreateDeliveryExecutiveResponse{
				DeliveryExecutive: &pb.Deliver_Executive{
					DeliveryExecutiveId: 1,
					Name:                "John Doe",
					Location:            "New York",
					Password:            "password123",
					DeliveryOrders:      []*pb.Delivery_Order{},
				},
			},
			wantErr:  false,
			errorMsg: "",
		},
		// Add more test cases for failure scenarios
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rows()

			server := &server{DB: gormDB}

			got, err := server.CreateDeliveryExecutive(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDeliveryExecutive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errorMsg {
				t.Errorf("CreateDeliveryExecutive() error message = %v, expected %v", err.Error(), tt.errorMsg)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAssignOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error creating mock db: %v", err)
	defer db.Close()

	dialect := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialect, &gorm.Config{})
	assert.Nil(t, err, "Error creating mock gorm db: %v", err)

	tests := []struct {
		name     string
		req      *pb.AssignOrderRequest
		rows     func()
		want     *pb.AssignOrderResponse
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Assign order - Success",
			req: &pb.AssignOrderRequest{
				Order: &pb.Order{
					OrderId:             1,
					PickUpLocation:      "New York",
					Status:              "PICKED_UP",
					DeliveryExecutiveId: 1,
				},
			},
			rows: func() {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "location", "password", "availability"}).AddRow(1, "John Doe", "New York", "password123", true))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			want: &pb.AssignOrderResponse{
				Order: &pb.Order{
					OrderId:             1,
					PickUpLocation:      "New York",
					Status:              "ASSIGNED",
					DeliveryExecutiveId: 1,
				},
			},
			wantErr:  false,
			errorMsg: "",
		},
		// Add more test cases for failure scenarios
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rows()

			server := &server{DB: gormDB}

			got, err := server.AssignOrder(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("AssignOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errorMsg {
				t.Errorf("AssignOrder() error message = %v, expected %v", err.Error(), tt.errorMsg)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}


func TestUpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "Error creating mock db: %v", err)
	defer db.Close()

	dialect := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialect, &gorm.Config{})
	assert.Nil(t, err, "Error creating mock gorm db: %v", err)

	type args struct {
        ctx context.Context
        req *pb.UpdateStatusRequest
    }

	tests := []struct {
		name     string
		args     args
		rows     func()
		want     *pb.UpdateStatusResponse
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Update status - Success",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
					"authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte("John Doe:password123"))},
				}),
				req: &pb.UpdateStatusRequest{
					OrderId: 1,
					Status: "PICKED_UP",
				},
			},
			rows: func() {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "location", "password", "availability"}).AddRow(1, "John Doe", "New York", "password123", false))
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"order_id", "delivery_status", "delivery_executive_id"}).AddRow(1, "ASSIGNED", 1))
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			want: &pb.UpdateStatusResponse{
				OrderId: 1,
				Status:  "PICKED_UP",
			},
			wantErr:  false,
			errorMsg: "",
		},
		// Add more test cases for failure scenarios
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rows()

			server := &server{DB: gormDB}

			got, err := server.UpdateStatus(tt.args.ctx, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errorMsg {
				t.Errorf("UpdateStatus() error message = %v, expected %v", err.Error(), tt.errorMsg)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}