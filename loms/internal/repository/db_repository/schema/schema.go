package schema

import (
	"github.com/jackc/pgtype"
)

type Order struct {
	OrderID   pgtype.Int8      `db:"order_id"`
	Status    pgtype.Varchar   `db:"status"`
	UserID    pgtype.Int8      `db:"user_id"`
	CreatedAt pgtype.Timestamp `db:"created_at"`
	UpdatedAt pgtype.Timestamp `db:"updated_at"`
}

type OrderItem struct {
	Sku     pgtype.Int8    `db:"sku"`
	OrderID pgtype.Int8    `db:"order_id"`
	Count   pgtype.Int4    `db:"count"`
	Price   pgtype.Numeric `db:"price"`
}

type Warehouse struct {
	WarehouseID pgtype.Int4 `db:"warehouse_id"`
}

type WarehouseStocks struct {
	Sku         pgtype.Int8 `db:"sku"`
	WarehouseID pgtype.Int4 `db:"warehouse_id"`
	Count       pgtype.Int4 `db:"count"`
}

type WarehouseReservations struct {
	Sku         pgtype.Int8      `db:"sku"`
	WarehouseID pgtype.Int4      `db:"warehouse_id"`
	OrderID     pgtype.Int8      `db:"order_id"`
	Count       pgtype.Int4      `db:"count"`
	ExpiredAt   pgtype.Timestamp `db:"expired_at"`
}

type WarehouseStocksWithReservations struct {
	WarehouseID pgtype.Int4 `db:"warehouse_id"`
	Count       pgtype.Int4 `db:"count"`
	Reservation pgtype.Int8 `db:"reservation"`
}
