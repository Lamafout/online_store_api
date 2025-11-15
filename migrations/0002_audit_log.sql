-- +goose Up
CREATE TABLE IF NOT EXISTS audit_log_order (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    order_item_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    order_status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE v1_audit_log_order AS (
    id BIGINT,
    order_id BIGINT,
    order_item_id BIGINT,
    customer_id BIGINT,
    order_status TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
-- +goose Down
DROP TABLE IF EXISTS audit_log_order;
DROP TYPE IF EXISTS v1_audit_log_order;