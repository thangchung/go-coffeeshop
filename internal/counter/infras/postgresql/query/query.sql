-- name: GetAll :many

SELECT
    o.id,
    order_source,
    loyalty_member_id,
    order_status,
    l.id as "line_item_id",
    item_type,
    name,
    price,
    item_status,
    is_barista_order
FROM "order".orders o
    LEFT JOIN "order".line_items l ON o.id = l.order_id;

-- name: GetByID :many

SELECT
    o.id,
    order_source,
    loyalty_member_id,
    order_status,
    l.id as "line_item_id",
    item_type,
    name,
    price,
    item_status,
    is_barista_order
FROM "order".orders o
    LEFT JOIN "order".line_items l ON o.id = l.order_id
WHERE o.id = $1;

-- name: CreateOrder :one

INSERT INTO
    "order".orders (
        id,
        order_source,
        loyalty_member_id,
        order_status,
        updated
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: InsertItemLine :one

INSERT INTO
    "order".line_items (
        id,
        item_type,
        name,
        price,
        item_status,
        is_barista_order,
        order_id,
        created,
        updated
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: UpdateOrder :exec

UPDATE "order".orders
SET
    order_status = $2,
    updated = $3
WHERE id = $1;

-- name: UpdateItemLine :exec

UPDATE "order".line_items
SET
    item_status = $2,
    updated = $3
WHERE id = $1;