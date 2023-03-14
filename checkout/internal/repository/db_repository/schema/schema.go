package schema

import "github.com/jackc/pgtype"

type Cart struct {
	UserID pgtype.Int8 `db:"user_id"`
	Sku    pgtype.Int8 `db:"sku"`
	Count  pgtype.Int4 `db:"count"`
}
