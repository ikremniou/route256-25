// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: stock.sql

package query

import (
	"context"
)

const getStocksBySkuId = `-- name: GetStocksBySkuId :one
select sku,
    total_count,
    reserved
from stocks
where sku = $1
`

type GetStocksBySkuIdRow struct {
	Sku        int64
	TotalCount int64
	Reserved   int64
}

func (q *Queries) GetStocksBySkuId(ctx context.Context, sku int64) (GetStocksBySkuIdRow, error) {
	row := q.db.QueryRow(ctx, getStocksBySkuId, sku)
	var i GetStocksBySkuIdRow
	err := row.Scan(&i.Sku, &i.TotalCount, &i.Reserved)
	return i, err
}
