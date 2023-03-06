package model

type OrderItem struct {
	Sku   uint32
	Count uint32
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}
