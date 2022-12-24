START TRANSACTION;

-- DO $$ BEGIN IF NOT EXISTS(

--         SELECT 1

--         FROM pg_namespace

--         WHERE

--             nspname = 'barista'

--     ) THEN CREATE SCHEMA barista;

-- END IF;

-- END $$;

CREATE SCHEMA IF NOT EXISTS "barista";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    barista.barista_orders (
        id uuid NOT NULL DEFAULT (uuid_generate_v4()),
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
            CONSTRAINT pk_barista_orders PRIMARY KEY (id)
    );

CREATE UNIQUE INDEX ix_barista_orders_id ON barista.barista_orders (id);

COMMIT;