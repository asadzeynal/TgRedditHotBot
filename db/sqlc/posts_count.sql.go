// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: posts_count.sql

package db

import (
	"context"
)

const getTotalCount = `-- name: GetTotalCount :one
SELECT count FROM posts_count
`

func (q *Queries) GetTotalCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getTotalCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}
