package model

const (
	OrderStatusNew             = "new"
	OrderStatusAwaitingPayment = "awaiting payment"
	OrderStatusFailed          = "failed"
	OrderStatusPayed           = "payed"
	OrderStatusCancelled       = "cancelled"
	// internal statuses. !IMPORTANT if status ends with 'ing'
	// it is a transition status and is not reported to outbox!
	OrderStatusReserving  = "reserving"
	OrderStatusPaying     = "paying"
	OrderStatusCancelling = "cancelling"
)

type OrderModel struct {
	Id     int64
	Status string
	UserId int64
	Items  []OrderItem
}
