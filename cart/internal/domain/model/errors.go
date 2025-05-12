package model

import (
	"errors"
)

var ErrTheCartIsEmpty = errors.New("the cart is empty")
var ErrProductDoesNotExist = errors.New("product does not exist")
var ErrNotEnoughItemsInStock = errors.New("not enough items in stock")
var ErrorUserIdLessThanZero = errors.New("user id must be greater than 0")
var ErrorSkuIdLessThanZero = errors.New("sku id must be greater than 0")
var ErrTotalCountExceeded = errors.New("the count of items exceeds the maximum allowed")
var ErrCreateOrderPreconditionFailed = errors.New("failed to create order, precondition failed")
