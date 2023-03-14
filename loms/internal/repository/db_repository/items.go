package db_repository

import (
	"context"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/db_repository/schema"
	"route256/loms/internal/repository/db_repository/transactor"

	sq "github.com/Masterminds/squirrel"
)

type OrderItemsRepository struct {
	transactor.QueryEngineProvider
}

func NewOrderItemsRepository(provider transactor.QueryEngineProvider) *OrderItemsRepository {
	return &OrderItemsRepository{
		provider,
	}
}

const ordersItemsTable = "orders_items"

func (r *OrderItemsRepository) CreateItem(ctx context.Context, item model.OrderItem) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Insert(ordersItemsTable).
		Columns(
			"sku",
			"order_id",
			"count",
			"price",
		).
		Values(item.Sku, item.OrderID, item.Count, item.Price).
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

func (r *OrderItemsRepository) GetItemsByOrderId(ctx context.Context, orderId int64) ([]model.OrderItem, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Select(
			"sku",
			"count",
			"price::numeric",
		).
		From(ordersItemsTable).
		Where(sq.Eq{"order_id": orderId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var item schema.OrderItem
	var result []model.OrderItem

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&item.Sku, &item.Count, &item.Price)
		if err != nil {
			return nil, err
		}
		result = append(result, model.OrderItem{
			Sku:     uint32(item.Sku.Int),
			OrderID: orderId,
			Count:   uint32(item.Count.Int),
			Price:   uint32(item.Price.Exp),
		})
	}

	return result, nil
}
