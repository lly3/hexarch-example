package rpc

import (
	"context"
	"fmt"
	"hex/internal/adapters/grpc/pb"
	"hex/internal/adapters/repository/order_repository"
	"hex/internal/application/use_case"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

var mockCreateOrderRequestData = &pb.CreateOrderRequestData{
	OrderId:      "0",
	ClientId:     "0",
	RestaurantId: "0",
	State:        "waiting-for-payment",
}

var mockCreateOrderClient = &pb.CreateOrderClient{
	ClientId:   "0",
	ClientName: "John",
}

var mockCreateOrderRestaurant = &pb.CreateOrderRestaurant{
	RestaurantId:      "0",
	RestaurantName:    "Sanwa",
	RestaurantIsOpen:  "true",
	RestaurantOwnerId: "0",
}

var mockCreateOrderFoods = []*pb.CreateOrderFoods{
	{
		FoodId:   "0",
		FoodName: "Kowkaijaow",
		Price:    30,
		Quantity: 2,
	},
}

func init() {
	var err error
	lis = bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("failed to initiate dbase connection: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("Error ping mongo client", err.Error())
	}
	fmt.Println("Mongo client ping")

	orderRepo := order_repository.NewOrderRepository(client.Database("hex"))

	uc := use_case.New(orderRepo)

	gRPCAdapter := New(uc)

	pb.RegisterOrderServer(grpcServer, gRPCAdapter)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("test server start error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func getGRPCConnection(ctx context.Context, t *testing.T) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	return conn
}

func TestCreateOrder(t *testing.T) {
	ctx := context.Background()
	conn := getGRPCConnection(ctx, t)
	defer conn.Close()

	client := pb.NewOrderClient(conn)

	req := &pb.CreateOrderRequest{
		Order:      mockCreateOrderRequestData,
		Client:     mockCreateOrderClient,
		Restaurant: mockCreateOrderRestaurant,
		Foods:      mockCreateOrderFoods,
	}

	res, err := client.CreateOrder(ctx, req)
	if err != nil {
		t.Fatalf("expected: %v, got: %v", nil, err)
	}

	require.Equal(t, res.OrderId, "0")
}

func TestUpdateOrderStatus(t *testing.T) {
	ctx := context.Background()
	conn := getGRPCConnection(ctx, t)
	defer conn.Close()

	client := pb.NewOrderClient(conn)

	req := &pb.UpdateOrderRequest{
		OrderId: "0",
		State:   "waiting-for-confirmation",
	}

	_, err := client.UpdateOrderStatus(ctx, req)
	if err != nil {
		t.Fatalf("expected: %v, got: %v", nil, err)
	}
}
