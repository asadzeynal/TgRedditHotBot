// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: post_images.sql

package db

import (
	"context"
)

const createPostImage = `-- name: CreatePostImage :one
INSERT INTO post_images (
    post,
    url
) VALUES (
$1, $2
) RETURNING id, post, url
`

type CreatePostImageParams struct {
	Post string `json:"post"`
	Url  string `json:"url"`
}

func (q *Queries) CreatePostImage(ctx context.Context, arg CreatePostImageParams) (PostImage, error) {
	row := q.db.QueryRowContext(ctx, createPostImage, arg.Post, arg.Url)
	var i PostImage
	err := row.Scan(&i.ID, &i.Post, &i.Url)
	return i, err
}

const getImagesByPost = `-- name: GetImagesByPost :many
SELECT id, post, url FROM post_images
WHERE post = $1 LIMIT 10
`

func (q *Queries) GetImagesByPost(ctx context.Context, post string) ([]PostImage, error) {
	rows, err := q.db.QueryContext(ctx, getImagesByPost, post)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PostImage{}
	for rows.Next() {
		var i PostImage
		if err := rows.Scan(&i.ID, &i.Post, &i.Url); err != nil {
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