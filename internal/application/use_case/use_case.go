package use_case

import (
	"context"
	"hex/internal/application/domain"
)

type UseCase struct {
	orderRepository OrderRepository
}

type HexOrder struct {
	Client     domain.Client
	Restaurant domain.Restaurant
	Order      domain.Order
	Foods      []domain.Food
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order domain.Order) error
	UpdateOrderStatus(ctx context.Context, orderId string, order domain.Order) error
	GetOrderById(ctx context.Context, orderId string) (domain.Order, error)
}

func New(orderRepo OrderRepository) UseCase {
	return UseCase{
		orderRepository: orderRepo,
	}
}
