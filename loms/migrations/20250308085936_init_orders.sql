-- +goose Up
-- +goose StatementBegin
create table orders (
    id bigserial primary key,
    user_id bigint not null,
    status text not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

create table order_items (
    id bigserial primary key,
    order_id bigint not null,
    sku bigint not null,
    quantity bigint not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table order_items;
drop table orders;
-- +goose StatementEnd
