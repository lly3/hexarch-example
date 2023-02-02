package use_case

import (
	"context"
	"errors"
	"hex/internal/application/domain"
)

func (u UseCase) CreateOrder(ctx context.Context, client domain.Client, order HexOrder) (string, error) {
	// validation things
	if !order.Restaurant.IsOpen {
		return "", errors.New("The restaurant is not open now")
	}

	u.orderRepository.CreateOrder(ctx, order.Order)

	return order.Order.Id, nil
}
