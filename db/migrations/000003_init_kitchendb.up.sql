START TRANSACTION;

-- DO $$ BEGIN IF NOT EXISTS(

--         SELECT 1

--         FROM pg_namespace

--         WHERE

--             nspname = 'kitchen'

--     ) THEN CREATE SCHEMA "kitchen";

-- END IF;

-- END $$;

CREATE SCHEMA IF NOT EXISTS "kitchen";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    kitchen.kitchen_orders (
        id uuid NOT NULL DEFAULT (uuid_generate_v4()),
        order_id uuid NOT NULL,
        item_type integer NOT NULL,
        item_name text NOT NULL,
        time_up timestamp
        with
            time zone NOT NULL,
            created timestamp
        with
            time zone NOT NULL DEFAULT (now()),
            updated timestamp
        with
            time zone NULL,
            CONSTRAINT pk_kitchen_orders PRIMARY KEY (id)
    );

CREATE UNIQUE INDEX ix_kitchen_orders_id ON kitchen.kitchen_orders (id);

COMMIT;