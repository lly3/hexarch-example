package use_case

import (
	"context"
	"hex/internal/application/domain"
)

func (u UseCase) UpdateOrderStatus(ctx context.Context, client domain.Client, orderId string, state domain.State) error {
	// validation things
	order, err := u.orderRepository.GetOrderById(ctx, orderId)
	if err != nil {
		return err
	}

	if err := order.UpdateOrderStatus(state); err != nil {
		return err
	}

	u.orderRepository.UpdateOrderStatus(ctx, order.Id, order)

	return nil
}
