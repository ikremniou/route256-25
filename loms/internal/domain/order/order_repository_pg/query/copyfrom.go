// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: copyfrom.go

package query

import (
	"context"
)

// iteratorForCreateOrderItems implements pgx.CopyFromSource.
type iteratorForCreateOrderItems struct {
	rows                 []CreateOrderItemsParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateOrderItems) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateOrderItems) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].OrderID,
		r.rows[0].Sku,
		r.rows[0].Quantity,
	}, nil
}

func (r iteratorForCreateOrderItems) Err() error {
	return nil
}

func (q *Queries) CreateOrderItems(ctx context.Context, arg []CreateOrderItemsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"order_items"}, []string{"order_id", "sku", "quantity"}, &iteratorForCreateOrderItems{rows: arg})
}
