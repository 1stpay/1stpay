-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE network_type AS ENUM ('evm', 'tron', 'ton', 'solana');
CREATE TYPE payment_status AS ENUM ('pending', 'completed', 'failed', 'not_filled');
CREATE TYPE payment_aml_status AS ENUM ('passed', 'failed', 'pending');

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    email text NOT NULL UNIQUE,
    password text NOT NULL,
    role text NOT NULL
);

CREATE TABLE blockchains (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    logo text,
    is_active boolean NOT NULL DEFAULT false,
    chain_type network_type NOT NULL,
    config jsonb
);

CREATE TABLE tokens (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    symbol text NOT NULL,
    logo text,
    blockchain_id uuid NOT NULL REFERENCES blockchains(id),
    is_native boolean NOT NULL DEFAULT false,
    is_active boolean NOT NULL DEFAULT false,
    config jsonb
);

CREATE TABLE merchants (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    user_id uuid NOT NULL REFERENCES users(id),
    name text NOT NULL,
    commission_rate numeric(5,2) NOT NULL DEFAULT 0
);

CREATE TABLE merchant_tokens (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant_id uuid NOT NULL REFERENCES merchants(id),
    token_id uuid NOT NULL REFERENCES tokens(id),
    is_active boolean NOT NULL DEFAULT false,
    balance numeric(20,8) NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE payments (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    merchant_id uuid NOT NULL REFERENCES merchants(id),
    used_token_id uuid REFERENCES tokens(id),
    requested_amount numeric(20,8) NOT NULL DEFAULT 0,
    paid_amount numeric(20,8) NOT NULL DEFAULT 0,
    commission_amount numeric(20,8) NOT NULL DEFAULT 0,
    expires_at timestamptz,
    aml_status payment_aml_status,
    status payment_status NOT NULL DEFAULT 'pending',
    invoice_email text
);

CREATE TABLE payment_addresses (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    payment_id uuid NOT NULL REFERENCES payments(id),
    token_id uuid NOT NULL REFERENCES tokens(id),
    requested_amount_wei bigint NOT NULL DEFAULT 0,
    paid_amount_wei bigint NOT NULL DEFAULT 0,
    requested_amount numeric(20,8) NOT NULL DEFAULT 0,
    paid_amount numeric(20,8) NOT NULL DEFAULT 0,
    public_key text NOT NULL,
    private_key text NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS payment_addresses;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS merchant_tokens;
DROP TABLE IF EXISTS merchant_blockchains;
DROP TABLE IF EXISTS merchants;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS blockchains;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS payment_aml_status;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS network_type;