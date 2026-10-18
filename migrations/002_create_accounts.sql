CREATE TABLE IF NOT EXISTS accounts (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL REFERENCES users(id),
    account_number TEXT NOT NULL UNIQUE,
    currency       TEXT NOT NULL DEFAULT 'INR',
    balance        BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);