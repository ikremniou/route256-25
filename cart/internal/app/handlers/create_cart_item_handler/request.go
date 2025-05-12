package create_cart_item_handler

type CreateCartItemRequest struct {
	Count uint32 `json:"count" validate:"required,min=1"`
}
