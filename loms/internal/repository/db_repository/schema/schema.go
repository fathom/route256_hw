package schema

import (
	"github.com/jackc/pgtype"
)

type Order struct {
	OrderID   pgtype.Int8      `db:"order_id"`
	Status    string           `db:"status"`
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

type WarehouseStocks struct {
	Sku         pgtype.Int8 `db:"sku"`
	WarehouseID pgtype.Int4 `db:"warehouse_id"`
	Count       pgtype.Int4 `db:"count"`
}
