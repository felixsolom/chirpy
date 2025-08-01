// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: retrieveAllChirps.sql

package database

import (
	"context"
)

const retrieveAllChirps = `-- name: RetrieveAllChirps :many
SELECT id, created_at, updated_at, body, user_id FROM chirps
ORDER BY created_at
`

func (q *Queries) RetrieveAllChirps(ctx context.Context) ([]Chirp, error) {
	rows, err := q.db.QueryContext(ctx, retrieveAllChirps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Chirp
	for rows.Next() {
		var i Chirp
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Body,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
