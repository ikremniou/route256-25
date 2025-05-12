package model

import "fmt"

type ErrOrderNotFound struct {
	OrderId int64
}

func (e *ErrOrderNotFound) Error() string {
	return fmt.Sprintf("orderId is not found, having: %d, ", e.OrderId)
}

type ErrInvalidOrderId struct {
	OrderId int64
}

func (e *ErrInvalidOrderId) Error() string {
	return fmt.Sprintf("invalid orderId: %d, must be greater than 0", e.OrderId)
}

type ErrReservedStockFailed struct {
	OrderId int64
}

func (e *ErrReservedStockFailed) Error() string {
	return fmt.Sprintf("failed to reserve stock for order: %d", e.OrderId)
}

type ErrStockNotFound struct {
	Sku int64
}

func (e *ErrStockNotFound) Error() string {
	return fmt.Sprintf("stock not found, sku: %d", e.Sku)
}

type ErrStockOutOfBounds struct {
	Sku        int64
	Reserved   uint32
	TotalCount uint32
	Change     uint32
}

func (e *ErrStockOutOfBounds) Error() string {
	return fmt.Sprintf("charging %d stocks for sku %d overflows. reserved: %d, total: %d",
		e.Change, e.Sku, e.Reserved, e.TotalCount)
}

type ErrOrderItemOutOfBounds struct {
	OrderId int64
	Sku     int64
	Total   uint64
}

func (e *ErrOrderItemOutOfBounds) Error() string {
	return fmt.Sprintf("order item sku: %d, total: %d, overflows order: %d",
		e.Sku, e.Total, e.OrderId)
}

type ErrInvalidOrderStatus struct {
	OrderId int64
	Status  string
}

func (e *ErrInvalidOrderStatus) Error() string {
	return fmt.Sprintf("status not valid for the operation, order: %d, status: %s",
		e.OrderId, e.Status)
}

type ErrOrderStatusMismatch struct {
	OrderId       int64
	CurrentStatus string
	ExpectedState string
}

func (e *ErrOrderStatusMismatch) Error() string {
	return fmt.Sprintf("status mismatch for order: %d, current: %s, expected: %s",
		e.OrderId, e.CurrentStatus, e.ExpectedState)
}
