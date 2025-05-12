package get_cart_items_handler

type GetCartItemResponseItem struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
}

type GetCartItemsResponse struct {
	Items      []GetCartItemResponseItem `json:"items"`
	TotalPrice uint32                    `json:"total_price"`
}
