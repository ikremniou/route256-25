package model

type CartItemModel struct {
	UserId int64  `validate:"required,min=1"`
	SkuId  int64  `validate:"required,min=1"`
	Count  uint32 `validate:"required,min=1"`
}
