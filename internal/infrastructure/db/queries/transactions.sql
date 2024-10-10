-- name: CreateTransaction :one
INSERT INTO transactions (from_wallet_id, to_wallet_id, amount, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = $1 LIMIT 1;

-- name: UpdateTransactionStatus :one
UPDATE transactions
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;