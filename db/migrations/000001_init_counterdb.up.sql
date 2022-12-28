START TRANSACTION;

-- DO $$ BEGIN IF NOT EXISTS(

--         SELECT 1

--         FROM pg_namespace

--         WHERE

--             nspname = 'order'

--     ) THEN CREATE SCHEMA "order";

-- END IF;

-- END $$;

CREATE SCHEMA IF NOT EXISTS "order";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    "order".orders (
        id uuid NOT NULL DEFAULT (uuid_generate_v4()),
        order_source integer NOT NULL,
        loyalty_member_id uuid NOT NULL,
        order_status integer NOT NULL,
        updated timestamp
        with
            time zone NULL,
            CONSTRAINT pk_orders PRIMARY KEY (id)
    );

CREATE TABLE
    "order".line_items (
        id uuid NOT NULL DEFAULT (uuid_generate_v4()),
        item_type integer NOT NULL,
        name text NOT NULL,
        price numeric NOT NULL,
        item_status integer NOT NULL,
        is_barista_order boolean NOT NULL,
        order_id uuid NULL,
        created timestamp
        with
            time zone NOT NULL DEFAULT (now()),
            updated timestamp
        with
            time zone NULL,
            CONSTRAINT pk_line_items PRIMARY KEY (id),
            CONSTRAINT fk_line_items_orders_order_temp_id FOREIGN KEY (order_id) REFERENCES "order".orders (id)
    );

CREATE UNIQUE INDEX ix_line_items_id ON "order".line_items (id);

CREATE INDEX ix_line_items_order_id ON "order".line_items (order_id);

CREATE UNIQUE INDEX ix_orders_id ON "order".orders (id);

COMMIT;