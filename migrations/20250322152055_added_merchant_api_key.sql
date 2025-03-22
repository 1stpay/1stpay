-- +goose Up
CREATE TABLE merchant_api_keys (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant_id uuid NOT NULL REFERENCES merchants(id),
    name text NOT NULL,
    api_key text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz,
    is_active boolean NOT NULL DEFAULT true
);

-- +goose Down
DROP TABLE IF EXISTS merchant_api_keys;