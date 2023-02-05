package order_repository

import (
	"context"
	"hex/internal/application/domain"
	"hex/internal/application/use_case"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodbOrder struct {
	Id           string `bson: id`
	ClientId     string `bson: client_id`
	RestaurantId string `bson: restaurant_id`
	State        string `bson: state`
}

type mongoDB struct {
	col *mongo.Collection
}

func newMongodbOrder(order domain.Order) mongodbOrder {
	return mongodbOrder{
		Id:           order.Id,
		ClientId:     order.ClientId,
		RestaurantId: order.RestaurantId,
		State:        string(order.State),
	}
}

func (m mongoDB) CreateOrder(ctx context.Context, order domain.Order) error {
	_, err := m.col.InsertOne(ctx, newMongodbOrder(order))
	if err != nil {
		return err
	}

	return nil
}

func (m mongoDB) UpdateOrderStatus(ctx context.Context, orderId string, order domain.Order) error {
	_, err := m.col.UpdateOne(ctx, bson.M{"id": orderId},
		bson.D{
			{"$set", bson.D{{"state", string(order.State)}}},
		})
	if err != nil {
		return err
	}

	return nil
}

func (m mongoDB) GetOrderById(ctx context.Context, orderId string) (domain.Order, error) {
	var order domain.Order
	err := m.col.FindOne(ctx, bson.M{"id": orderId}).Decode(&order)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func NewOrderRepository(db *mongo.Database) use_case.OrderRepository {
	return &mongoDB{
		col: db.Collection("orders"),
	}
}
