-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    delivery_address TEXT NOT NULL,
    total_price_cents BIGINT NOT NULL,
    total_price_currency TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_order_customer_id ON orders (customer_id);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL,
    product_title TEXT NOT NULL,
    product_url TEXT NOT NULL,
    price_cents BIGINT NOT NULL,
    price_currency TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_order_item_order_id ON order_items (order_id);

CREATE TYPE v1_order AS (
    id BIGINT,
    customer_id BIGINT,
    delivery_address TEXT,
    total_price_cents BIGINT,
    total_price_currency TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TYPE v1_order_item AS (
    id BIGINT,
    order_id BIGINT,
    product_id BIGINT,
    quantity INTEGER,
    product_title TEXT,
    product_url TEXT,
    price_cents BIGINT,
    price_currency TEXT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS audit_log_order (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    order_item_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    order_status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- +goose Down
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS v1_order;
DROP TYPE IF EXISTS v1_order_item;