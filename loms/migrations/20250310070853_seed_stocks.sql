-- +goose Up
-- +goose StatementBegin

insert into stocks (sku, total_count, reserved, created_at, updated_at)
values
    (139275865, 65534, 0, now(), now()),
    (1076963, 65534, 10, now(), now()),
    (1148162, 200, 20, now(), now()),
    (1625903, 250, 30, now(), now()),
    (2618151, 300, 40, now(), now()),
    (2956315, 350, 50, now(), now());

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

delete from stocks where sku in (139275865, 1076963, 1148162, 1625903, 2618151, 2956315);

-- +goose StatementEnd
