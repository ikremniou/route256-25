-- name: UpdateRowStatus :batchexec
update outbox
set status = @status, updated_at = now()
where id = @id;

-- name: GetPending :many
select * from outbox
where status = 'pending'
order by updated_at asc
limit @batch_size
for update skip locked;
