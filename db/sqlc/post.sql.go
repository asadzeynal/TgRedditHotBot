// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: post.sql

package db

import (
	"context"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (
    id,
    title,
    url
) VALUES (
$1, $2, $3
) RETURNING id, title, url, created_at
`

type CreatePostParams struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRow(ctx, createPost, arg.ID, arg.Title, arg.Url)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Url,
		&i.CreatedAt,
	)
	return i, err
}

const getRandomPost = `-- name: GetRandomPost :one
SELECT id, title, url, created_at FROM posts OFFSET $1 LIMIT 1
`

func (q *Queries) GetRandomPost(ctx context.Context, offset int32) (Post, error) {
	row := q.db.QueryRow(ctx, getRandomPost, offset)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Url,
		&i.CreatedAt,
	)
	return i, err
}
