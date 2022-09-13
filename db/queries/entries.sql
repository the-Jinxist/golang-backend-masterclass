-- name: CreateEntry :one
INSERT INTO entries (
  account_id, amount
) VALUES (
    $1, $2
)
RETURNING *;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY created_at
LIMIT $2
OFFSET $3;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1
LIMIT 1;