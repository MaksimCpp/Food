package domain

const (
    StatusAccepted = "accepted"
    StatusOutForDelivery = "out_for_delivery"
    StatusDelivered = "delivered"
)

type Order struct {
	UserID           int64
    ProductID        int64
    ProductsQuantity int64
    CourierID        int64
    Address          string
    Status           string
}