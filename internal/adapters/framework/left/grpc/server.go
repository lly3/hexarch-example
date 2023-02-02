package rpc

import (
	"fmt"
	"hex/internal/adapters/framework/left/grpc/pb"
	"hex/internal/ports"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Adapter struct {
	useCase ports.UseCase
}

func New(useCase ports.UseCase) *Adapter {
	return &Adapter{
		useCase: useCase,
	}
}

func (grpca Adapter) Run() {

	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to lesten on port 9000: %v", err)
	}

	orderServiceServer := grpca
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServer(grpcServer, orderServiceServer)

	fmt.Println("server is running on port 9000")

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc server on port 9000: %v", err)
	}
}
