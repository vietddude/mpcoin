-- name: CreateWallet :one
INSERT INTO wallets (user_id, address, encrypted_private_key)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetWallet :one
SELECT * FROM wallets
WHERE id = $1 LIMIT 1;

-- name: GetWalletByAddress :one
SELECT * FROM wallets
WHERE address = $1 LIMIT 1;

-- name: GetWalletByUserID :one
SELECT * FROM wallets
WHERE user_id = $1 LIMIT 1;
