-- +goose Up
-- +goose StatementBegin
alter table stocks
add constraint stocks_totalCount_uint32 check (
        total_count >= 0
        and total_count <= 4294967295
    ),
    add constraint stocks_reserved_uint32 check (
        reserved >= 0
        and reserved <= 4294967295
    ),
    add constraint stocks_totalCount_gte_reserved check (total_count >= reserved);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table stocks drop constraint stocks_totalCount_uint32,
    drop constraint stocks_reserved_uint32,
    drop constraint stocks_totalCount_gte_reserved;
-- +goose StatementEnd