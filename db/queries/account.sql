-- name: CreateAccount :one
INSERT INTO accounts (
  owner, balance, currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY created_at;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;