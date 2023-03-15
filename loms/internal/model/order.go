package model

import "time"

type Order struct {
	OrderID   int64
	Status    OrderStatus
	UserID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderItem struct {
	Sku     uint32
	OrderID int64
	Count   uint32
	Price   uint32
}

type OrderStatus string

const (
	NewStatus       OrderStatus = `new`
	AwaitingPayment OrderStatus = `awaiting payment`
	Failed          OrderStatus = `failed`
	Payed           OrderStatus = `payed`
	Cancelled       OrderStatus = `cancelled`
)
