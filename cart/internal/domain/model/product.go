package model

type ProductModel struct {
	SkuId int64  `json:"sku"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}
