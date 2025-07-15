-- name: UserIdFromEmail :one 
SELECT id FROM users
WHERE email = $1;