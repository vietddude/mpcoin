-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create chains table
CREATE TABLE chains (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    chain_id VARCHAR(255) UNIQUE NOT NULL,
    rpc_url VARCHAR(255) NOT NULL,
    native_currency VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create wallets table
CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    address VARCHAR(42) NOT NULL,
    encrypted_private_key BYTEA NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
);

-- Create tokens table
CREATE TABLE tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chain_id UUID NOT NULL,
    contract_address VARCHAR(42) NOT NULL,
    name VARCHAR(255) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    decimals INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_chain
        FOREIGN KEY (chain_id)
        REFERENCES chains (id)
        ON DELETE CASCADE
);

-- Create balances table
CREATE TABLE balances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID NOT NULL,
    chain_id UUID NOT NULL,
    token_id UUID,
    balance NUMERIC(78, 0) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_wallet
        FOREIGN KEY (wallet_id)
        REFERENCES wallets (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_token
        FOREIGN KEY (token_id)
        REFERENCES tokens (id)
        ON DELETE SET NULL,
    CONSTRAINT fk_chain_balance
        FOREIGN KEY (chain_id)
        REFERENCES chains (id)
        ON DELETE CASCADE
);

-- Create transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID NOT NULL,
    chain_id UUID NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42) NOT NULL,
    amount NUMERIC(78, 0) NOT NULL,
    token_id UUID,
    gas_price NUMERIC(78, 0),
    gas_limit BIGINT,
    nonce BIGINT,
    status VARCHAR(20) NOT NULL,
    tx_hash VARCHAR(66),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_wallet_transaction
        FOREIGN KEY (wallet_id)
        REFERENCES wallets (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_chain_transaction
        FOREIGN KEY (chain_id)
        REFERENCES chains (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_token_transaction
        FOREIGN KEY (token_id)
        REFERENCES tokens (id)
        ON DELETE SET NULL
);

-- Insert chains
INSERT INTO chains (id, name, chain_id, rpc_url, native_currency) VALUES
(uuid_generate_v4(), 'Sepolia', '11511', 'https://eth-sepolia.g.alchemy.com/v2/demo', 'ETH'),
(uuid_generate_v4(), 'Base Sepolia', '84532', 'https://sepolia.base.org', 'ETH'),
(uuid_generate_v4(), 'Arbitrum Sepolia', '421613', 'https://sepolia.arbitrum.io/rpc', 'ETH');

-- Insert tokens
INSERT INTO tokens (id, chain_id, contract_address, name, symbol, decimals)
SELECT 
    uuid_generate_v4(),
    chains.id,
    '0x1',
    'Ethereum',
    'ETH',
    18
FROM chains
WHERE chains.chain_id IN ('11511', '84532', '421613');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS balances;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS chains;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
