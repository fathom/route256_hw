package domain

type Domain struct {
}

func New() *Domain {
	return &Domain{}
}

type OrderItem struct {
	Sku   uint32
	Count uint32
}

type OrderStatus string

const (
	NewStatus       OrderStatus = `new`
	AwaitingPayment OrderStatus = `awaiting payment`
	Failed          OrderStatus = `failed`
	Payed           OrderStatus = `payed`
	Cancelled       OrderStatus = `cancelled`
)
