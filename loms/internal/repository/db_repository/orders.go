package db_repository

import (
	"context"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/db_repository/schema"
	"route256/loms/internal/repository/db_repository/transactor"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type OrdersRepository struct {
	transactor.QueryEngineProvider
}

func NewOrdersRepository(provider transactor.QueryEngineProvider) *OrdersRepository {
	return &OrdersRepository{
		provider,
	}
}

const ordersTable = "orders"

func (r *OrdersRepository) GetOrder(ctx context.Context, orderId int64) (model.Order, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Select(
			"status",
			"user_id",
			"created_at",
			"updated_at").
		From(ordersTable).
		Where(sq.Eq{"order_id": orderId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return model.Order{}, err
	}

	var order schema.Order
	err = db.QueryRow(ctx, query, args...).Scan(&order.Status, &order.UserID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return model.Order{}, err
	}

	return model.Order{
		OrderID:   orderId,
		Status:    model.OrderStatus(order.Status.String),
		UserID:    order.UserID.Int,
		CreatedAt: order.CreatedAt.Time,
		UpdatedAt: order.UpdatedAt.Time,
	}, nil
}

func (r *OrdersRepository) CreateOrder(ctx context.Context, order model.Order) (int64, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Insert(ordersTable).
		Columns(
			"status",
			"user_id",
			"created_at",
			"updated_at").
		Values(order.Status, order.UserID, time.Now(), time.Now()).
		Suffix("RETURNING order_id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, err
	}

	var orderID int64
	err = db.QueryRow(ctx, query, args...).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (r *OrdersRepository) UpdateStatusOrder(ctx context.Context, orderId int64, status model.OrderStatus) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Update(ordersTable).
		Set("status", string(status)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"order_id": orderId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
