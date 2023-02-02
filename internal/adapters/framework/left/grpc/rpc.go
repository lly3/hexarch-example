package rpc

import (
	"context"
	"hex/internal/adapters/framework/left/grpc/pb"
	"hex/internal/application/domain"
	"hex/internal/application/use_case"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func toHexOrderUseCase(res *pb.CreateOrderRequest) use_case.HexOrder {
	var foods []domain.Food

	for _, pb_food := range res.Foods {
		food := domain.Food{
			Id:       pb_food.FoodId,
			Name:     pb_food.FoodName,
			Price:    pb_food.Price,
			Quantity: int(pb_food.Quantity),
		}
		foods = append(foods, food)
	}

	return use_case.HexOrder{
		Client: domain.Client{
			Id:   res.Client.ClientId,
			Name: res.Client.ClientName,
		},
		Restaurant: domain.Restaurant{
			Id:      res.Restaurant.RestaurantId,
			Name:    res.Restaurant.RestaurantName,
			IsOpen:  res.Restaurant.RestaurantIsOpen == "true",
			OwnerId: res.Restaurant.RestaurantOwnerId,
		},
		Order: domain.Order{
			Id:           res.Order.OrderId,
			ClientId:     res.Order.ClientId,
			RestaurantId: res.Order.RestaurantId,
			State:        domain.State(res.Order.State),
		},
		Foods: foods,
	}
}

func (grpca Adapter) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	res := &pb.CreateOrderResponse{}

	order := toHexOrderUseCase(req)

	orderId, err := grpca.useCase.CreateOrder(ctx, domain.Client{}, order)
	if err != nil {
		return res, status.Error(codes.Internal, "unexpected error")
	}

	res = &pb.CreateOrderResponse{
		OrderId: orderId,
	}

	return res, nil

}

func (grpca Adapter) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderRequest) (*emptypb.Empty, error) {
	res := &emptypb.Empty{}

	err := grpca.useCase.UpdateOrderStatus(ctx, domain.Client{}, req.OrderId, domain.State(req.State))
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}
