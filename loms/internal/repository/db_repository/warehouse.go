package db_repository

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"log"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/db_repository/schema"
	"route256/loms/internal/repository/db_repository/transactor"
	"time"
)

type WarehouseRepository struct {
	transactor.QueryEngineProvider
}

func NewWarehouseRepository(provider transactor.QueryEngineProvider) *WarehouseRepository {
	return &WarehouseRepository{
		provider,
	}
}

const warehouseTable = "warehouse"
const warehouseStocksTable = "warehouse_stocks"
const warehouseReservationsTable = "warehouse_reservations"

var (
	ErrNotFoundFailed = errors.New("not found")
)

func (r *WarehouseRepository) GetStocksBySku(ctx context.Context, sku uint32) ([]model.StockItem, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Select(
			"s.warehouse_id",
			"s.count",
			"SUM(r.count) reservation",
		).
		From(fmt.Sprintf("%v s", warehouseStocksTable)).
		LeftJoin(fmt.Sprintf("%v r ON s.sku = r.sku and s.warehouse_id = r.warehouse_id and r.expired_at > ?", warehouseReservationsTable), time.Now()).
		Where(sq.Eq{"s.sku": sku}).
		GroupBy("s.warehouse_id, s.count").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var item schema.WarehouseStocksWithReservations
	var result []model.StockItem

	log.Printf("GetStocksBySku Query: %v with %v", query, args)
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&item.WarehouseID, &item.Count, &item.Reservation)
		if err != nil {
			return nil, err
		}
		result = append(result, model.StockItem{
			WarehouseID: int64(item.WarehouseID.Int),
			Count:       uint32(item.Count.Int),
			Reservation: uint32(item.Reservation.Int),
		})
	}

	return result, nil
}

func (r *WarehouseRepository) GetStocks(ctx context.Context, sku uint32, warehouseId int64) (model.StockItem, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Select(
			"sku",
			"warehouse_id",
			"count",
		).
		From(warehouseStocksTable).
		Where(sq.Eq{"sku": sku}).
		Where(sq.Eq{"warehouse_id": warehouseId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return model.StockItem{}, err
	}

	log.Printf("GetStocks Query: %v with %v", query, args)
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return model.StockItem{}, err
	}
	defer rows.Close()

	var item schema.WarehouseStocks
	if rows.Next() {
		err = rows.Scan(&item.Sku, &item.WarehouseID, &item.Count)
		if err != nil {
			return model.StockItem{}, err
		}
		return model.StockItem{
			WarehouseID: int64(item.WarehouseID.Int),
			Count:       uint32(item.Count.Int),
		}, nil
	}

	return model.StockItem{}, ErrNotFoundFailed
}

func (r *WarehouseRepository) UpdateWarehouse(ctx context.Context, sku uint32, warehouseId int64, count uint32) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Update(warehouseStocksTable).
		Set("count", count).
		Where(sq.Eq{"sku": sku}).
		Where(sq.Eq{"warehouse_id": warehouseId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	log.Printf("UpdateWarehouse Query: %v with %v", query, args)
	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *WarehouseRepository) GetReservationByOrderId(ctx context.Context, orderId int64) ([]model.StockReservationItem, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Select(
			"sku",
			"warehouse_id",
			"count",
		).
		From(warehouseReservationsTable).
		Where(sq.Eq{"order_id": orderId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var item schema.WarehouseReservations
	var result []model.StockReservationItem

	log.Printf("GetReservationByOrderId Query: %v with %v", query, args)
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&item.Sku, &item.WarehouseID, &item.Count)
		if err != nil {
			return nil, err
		}
		result = append(result, model.StockReservationItem{
			Sku:         uint32(item.Sku.Int),
			WarehouseID: int64(item.WarehouseID.Int),
			Count:       uint32(item.Count.Int),
		})
	}

	return result, nil
}

func (r *WarehouseRepository) AddReservation(ctx context.Context, sku uint32, warehouseId int64, orderId int64, count uint32, expiredAt time.Time) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Insert(warehouseReservationsTable).
		Columns(
			"sku",
			"warehouse_id",
			"order_id",
			"count",
			"expired_at",
		).
		Values(sku, warehouseId, orderId, count, expiredAt).
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

func (r *WarehouseRepository) DeleteReservation(ctx context.Context, orderId int64) error {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Delete(warehouseReservationsTable).
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
