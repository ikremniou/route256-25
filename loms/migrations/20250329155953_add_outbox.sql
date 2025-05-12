-- +goose Up
-- +goose StatementBegin
create table outbox (
    id bigserial primary key,
    aggregate_id text not null,
    aggregate_type text not null,
    event_type text not null,
    key text not null,
    payload jsonb not null,
    topic text not null,
    status text not null default 'pending',
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

create index idx_outbox_pending on outbox (status, updated_at)
where status = 'pending';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table outbox;
-- +goose StatementEnd
