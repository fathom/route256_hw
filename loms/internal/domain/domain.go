package domain

type Domain struct {
}

func New() *Domain {
	return &Domain{}
}

type OrderItem struct {
	Sku   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type OrderStatus string

const (
	NewStatus       OrderStatus = `new`
	AwaitingPayment OrderStatus = `awaiting payment`
	Failed          OrderStatus = `failed`
	Payed           OrderStatus = `payed`
	Cancelled       OrderStatus = `cancelled`
)
