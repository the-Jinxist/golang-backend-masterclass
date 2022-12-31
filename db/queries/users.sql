-- name: CreateUser :one
INSERT INTO users (
  username, hashed_password, full_name, email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1;

-- name: UpdateUser :one
-- This method uses optional parameters with those weird CASE WHEN ELSE parameters
-- UPDATE users 
-- set 
--   hashed_password = CASE
--     WHEN @set_hashed_password::boolean = TRUE THEN @hashed_password
--     ELSE hashed_password
--   END,
--   full_name = CASE
--     WHEN @set_full_name::boolean = TRUE THEN @full_name
--     ELSE full_name
--   END,
--   email = CASE
--     WHEN @set_email::boolean = TRUE THEN @email
--     ELSE email
--   END
-- WHERE 
--   username = @username
-- RETURNING *;


-- The next method uses nullable paramters method to create the same updateUser method with

-- name: UpdateUser :one
UPDATE users
SET
 hashed_password = coalesce(sqlc.narg('hashed_password'), hashed_password),
 full_name = coalesce(sqlc.narg('full_name'), full_name),
 email = coalesce(sqlc.narg('email'), email),
 password_changed_at = coalesce(sqlc.narg('password_changed_at'), password_changed_at)
WHERE username = sqlc.arg('username')
RETURNING *;