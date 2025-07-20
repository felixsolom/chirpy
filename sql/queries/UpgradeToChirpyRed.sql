-- name: UpgradeToChirpyRed :one
UPDATE users
SET 
    updated_at = NOW(),
    is_chirpy_red = TRUE
WHERE id=$1
RETURNING *; 

