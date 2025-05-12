-- name: GetStocksBySkuId :one
select sku,
    total_count,
    reserved
from stocks
where sku = $1;

-- name: Reserve :batchone
update stocks
set reserved = reserved + @count
where sku = $1
returning sku;

-- name: RemoveReserve :batchexec
update stocks
set reserved = reserved - @count, total_count = total_count - @count
where sku = $1;
