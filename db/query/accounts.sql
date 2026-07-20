-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2,$3
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1 
FOR NO KEY UPDATE;
/*The FOR UPDATE clause in SQL is used to lock the selected rows so that other transactions cannot modify or delete them until the current transaction is committed or rolled back. It is commonly used in transactional systems to prevent concurrency issues like lost updates.*/



-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1 
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
  set balance = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;
