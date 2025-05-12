-- +goose Up
-- +goose StatementBegin
create table stocks (
    id serial primary key,
    sku bigint not null,
    total_count bigint not null,
    reserved bigint not null,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table stocks;
-- +goose StatementEnd
