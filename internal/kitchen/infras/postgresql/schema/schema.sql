CREATE SCHEMA "kitchen";

CREATE TABLE
    kitchen.kitchen_orders (
        id uuid NOT NULL DEFAULT uuid_generate_v4(),
        order_id uuid NOT NULL,
        item_type int4 NOT NULL,
        item_name text NOT NULL,
        time_up timestamptz NOT NULL,
        created timestamptz NOT NULL DEFAULT now(),
        updated timestamptz NULL,
        CONSTRAINT pk_kitchen_orders PRIMARY KEY (id)
    );

CREATE UNIQUE INDEX ix_kitchen_orders_id ON kitchen.kitchen_orders USING btree (id);