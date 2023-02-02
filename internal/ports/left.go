package ports

import (
	"context"
	"hex/internal/application/domain"
	"hex/internal/application/use_case"
)

type UseCase interface {
	CreateOrder(ctx context.Context, client domain.Client, order use_case.HexOrder) (string, error)
	UpdateOrderStatus(ctx context.Context, client domain.Client, orderId string, status domain.State) error
}
