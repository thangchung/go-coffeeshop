-- name: CreateOrder :one

INSERT INTO
    barista.barista_orders (
        id,
        item_type,
        item_name,
        time_up,
        created,
        updated
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;