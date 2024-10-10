-- name: CreateWallet :one
INSERT INTO wallets (user_id, address, public_key, encrypted_private_key)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWallet :one
SELECT * FROM wallets
WHERE id = $1 LIMIT 1;

-- name: UpdateWalletBalance :one
UPDATE wallets
SET balance = balance + $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
