-- +goose Up
-- +goose StatementBegin
create table if not exists orders
(
    order_id   bigserial
        constraint orders_pk
            primary key,
    status     varchar,
    user_id    bigint,
    created_at timestamp,
    updated_at  timestamp
);

create index if not exists status_index
    on orders (status);

create index if not exists user_id_index
    on orders (user_id);


create table if not exists orders_items
(
    order_id bigint,
    sku      bigint,
    count    integer,
    price    money,
    constraint orders_items_pk
        primary key (order_id, sku)
);

create table if not exists warehouse
(
    warehouse_id serial
        constraint warehouse_pk
            primary key
);

create table if not exists warehouse_stocks
(
    sku          bigint,
    warehouse_id integer,
    count        integer,
    reserved     integer,
    constraint warehouse_stocks_pk
        primary key (sku, warehouse_id)
);

create index if not exists warehouse_id_index
    on warehouse (warehouse_id);

create table if not exists warehouse_reservations
(
    sku          bigint,
    warehouse_id integer,
    order_id     bigint,
    count        integer,
    expired_at   timestamp,
    constraint warehouse_reservations_pk
        primary key (sku, warehouse_id, order_id)
);

create index warehouse_reservations_expired_at_index
    on warehouse_reservations (expired_at);

create index warehouse_reservations_order_id_index
    on warehouse_reservations (order_id);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists orders;
drop table if exists orders_items;
drop table if exists warehouse;
drop table if exists warehouse_stocks;
drop table if exists warehouse_reservations;
-- +goose StatementEnd
