-- name: GetUserFromEmail :one
SELECT * FROM users
WHERE email=$1;