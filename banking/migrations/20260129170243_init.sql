-- +goose Up

CREATE TYPE currency AS ENUM ('USD', 'EUR');
CREATE TYPE transaction_type AS ENUM ('transfer', 'exchange');

CREATE TABLE users (
  id UUID PRIMARY KEY NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE accounts (
  id UUID PRIMARY KEY NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  currency currency NOT NULL,
  balance NUMERIC(19, 2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE transactions (
  id UUID PRIMARY KEY NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type transaction_type NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE ledger (
  id SERIAL PRIMARY KEY NOT NULL,
  account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
  transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
  amount NUMERIC(19, 2) NOT NULL,
  currency currency NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE UNIQUE INDEX idx_accounts_user_currency ON accounts(user_id, currency);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_ledger_account_id ON ledger(account_id);
CREATE INDEX idx_ledger_transaction_id ON ledger(transaction_id);
CREATE INDEX idx_ledger_created_at ON ledger(created_at);

-- +goose Down

DROP INDEX IF EXISTS idx_ledger_created_at;
DROP INDEX IF EXISTS idx_ledger_transaction_id;
DROP INDEX IF EXISTS idx_ledger_account_id;
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_user_id;
DROP INDEX IF EXISTS idx_accounts_user_id;

DROP TABLE IF EXISTS ledger;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS currency;