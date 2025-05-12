package model

type OrderItem struct {
	Sku   int64  `validate:"required,gte=1"`
	Count uint32 `validate:"required,gte=1"`
}

type CreateOrderModel struct {
	UserId int64       `validate:"required,gte=1"`
	Items  []OrderItem `validate:"required,min=1,dive"`
}
