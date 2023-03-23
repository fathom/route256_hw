package domain

import (
	"context"
	"route256/loms/internal/model"
	"time"
)

type OrdersRepository interface {
	GetOrder(context.Context, int64) (model.Order, error)
	CreateOrder(context.Context, model.Order) (int64, error)
	UpdateStatusOrder(context.Context, int64, model.OrderStatus) error
}

type OrderItemsRepository interface {
	CreateItem(context.Context, model.OrderItem) error
	GetItemsByOrderId(context.Context, int64) ([]model.OrderItem, error)
}

type WarehouseRepository interface {
	GetStocksBySku(ctx context.Context, sku uint32) ([]model.StockItem, error)
	UpdateWarehouse(ctx context.Context, sku uint32, warehouseId int64, count uint32) error
	AddReservation(ctx context.Context, sku uint32, warehouseId int64, orderId int64, count uint32, expiredAt time.Time) error
	DeleteReservation(ctx context.Context, orderId int64) error
	GetReservationByOrderId(ctx context.Context, orderId int64) ([]model.StockReservationItem, error)
	GetStocks(ctx context.Context, sku uint32, warehouseId int64) (model.StockItem, error)
}

type DeleteReservationWorker interface {
	AddDelayJob(job model.JobDeleteReservation)
}

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, f func(ctxTX context.Context) error) error
	// todo RunSerializable()
	// todo RunCommitted()
	// todo RunUncommitted()
}

type Domain struct {
	TransactionManager
	OrdersRepository
	OrderItemsRepository
	WarehouseRepository
	DeleteReservationWorker
}

func New(
	transactionManager TransactionManager,
	ordersRepository OrdersRepository,
	orderItemsRepository OrderItemsRepository,
	warehouseRepository WarehouseRepository,
	deleteReservationWorker DeleteReservationWorker,
) *Domain {
	return &Domain{
		transactionManager,
		ordersRepository,
		orderItemsRepository,
		warehouseRepository,
		deleteReservationWorker,
	}
}
