-- +goose Up
-- +goose StatementBegin
create table  if not exists cart
(
    user_id bigint,
    sku     bigint,
    count   integer,
    constraint cart_pk
        primary key (user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists orders;
-- +goose StatementEnd
