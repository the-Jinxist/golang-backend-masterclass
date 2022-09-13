-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account, to_account, amout
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetTransfers :many
SELECT * FROM transfers 
WHERE from_account = $1 OR to_account = $2
ORDER BY created_at
LIMIT $3
OFFSET $4;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1
LIMIT 1;