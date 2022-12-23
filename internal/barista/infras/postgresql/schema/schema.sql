CREATE SCHEMA "barista";

CREATE TABLE
    barista.barista_orders (
        id uuid NOT NULL DEFAULT uuid_generate_v4(),
        item_type int4 NOT NULL,
        item_name text NOT NULL,
        time_up timestamptz NOT NULL,
        created timestamptz NOT NULL DEFAULT now(),
        updated timestamptz NULL,
        CONSTRAINT pk_barista_orders PRIMARY KEY (id)
    );

CREATE UNIQUE INDEX ix_barista_orders_id ON barista.barista_orders USING btree (id);