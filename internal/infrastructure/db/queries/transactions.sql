-- name: CreateTransaction :one
INSERT INTO transactions (id, wallet_id , chain_id, to_address, amount, token_id, gas_price, gas_limit, nonce, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1 LIMIT 1;

-- name: GetTransactionsByWalletID :many
SELECT * FROM transactions
WHERE wallet_id = $1;

-- name: UpdateTransaction :one
UPDATE transactions 
SET (status, tx_hash, gas_price, gas_limit, nonce) = ($2, $3, $4, $5, $6)
WHERE id = $1
RETURNING *;