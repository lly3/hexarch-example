package main

import (
	"context"
	"fmt"
	rpc "hex/internal/adapters/framework/left/grpc"
	"hex/internal/adapters/framework/right/db/order_repository"
	"hex/internal/application/use_case"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
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

	grpcServer := rpc.New(uc)

	grpcServer.Run()
}
