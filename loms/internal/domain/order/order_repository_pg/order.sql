-- name: CreateOrder :one
with create_order as (
    insert into orders (user_id, status)
    values ($1, $2)
    returning id, status, user_id
), insert_outbox as (
    insert into outbox (
        aggregate_id,
        aggregate_type,
        event_type,
        key,
        payload,
        topic,
        status
    )
    select
        co.id,
        'order',
        'order_status_updated',
        co.user_id::text,
        json_build_object(
            'from_status', '',
            'to_status', co.status,
            'order_id', co.id,
            'user_id', co.user_id
        ),
        @topic,
        'pending'
    from create_order as co
)
select id from create_order;

-- name: CreateOrderItems :copyfrom
insert into order_items (order_id, sku, quantity)
values ($1, $2, $3);

-- name: GetByID :many
select o.id,
    o.user_id,
    o.status,
    oi.sku,
    oi.quantity
from orders as o
    join order_items as oi on o.id = oi.order_id
where o.id = $1
order by oi.sku asc;

-- name: UpdateStatus :one
with update_status as (
    update orders as o
    set status = @status,
        updated_at = now()
    where o.id = @id
        and o.status = @old_status
    returning id, status, user_id
), insert_outbox as (
    insert into outbox (
        aggregate_id,
        aggregate_type,
        event_type,
        key,
        payload,
        topic,
        status
    )
    select
        u.id,
        'order',
        'order_status_updated',
        u.user_id::text,
        json_build_object(
            'from_status', @old_status,
            'to_status', @status,
            'order_id', u.id,
            'user_id', u.user_id
        ),
        @topic,
        'pending'
    from update_status as u
    where u.status not like '%ing'
)
select id
from update_status;