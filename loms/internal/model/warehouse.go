package model

import "time"

type StockItem struct {
	WarehouseID int64
	Count       uint32
	Reservation uint32
}

type StockReservationItem struct {
	Sku         uint32
	WarehouseID int64
	OrderId     int64
	Count       uint32
	ExpiredAt   time.Time
}

type JobDeleteReservation struct {
	OrderId int64
}
