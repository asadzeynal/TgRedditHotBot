// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: post_images.sql

package db

import (
	"context"
)

const createPostImage = `-- name: CreatePostImage :one
INSERT INTO post_images (
    post,
    url,
    is_gif
) VALUES (
$1, $2, $3
) RETURNING id, post, url, is_gif, tg_file_id
`

type CreatePostImageParams struct {
	Post  string `json:"post"`
	Url   string `json:"url"`
	IsGif bool   `json:"is_gif"`
}

func (q *Queries) CreatePostImage(ctx context.Context, arg CreatePostImageParams) (PostImage, error) {
	row := q.db.QueryRow(ctx, createPostImage, arg.Post, arg.Url, arg.IsGif)
	var i PostImage
	err := row.Scan(
		&i.ID,
		&i.Post,
		&i.Url,
		&i.IsGif,
		&i.TgFileID,
	)
	return i, err
}

const getImagesByPost = `-- name: GetImagesByPost :many
SELECT id, post, url, is_gif, tg_file_id FROM post_images
WHERE post = $1 LIMIT 10
`

func (q *Queries) GetImagesByPost(ctx context.Context, post string) ([]PostImage, error) {
	rows, err := q.db.Query(ctx, getImagesByPost, post)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PostImage{}
	for rows.Next() {
		var i PostImage
		if err := rows.Scan(
			&i.ID,
			&i.Post,
			&i.Url,
			&i.IsGif,
			&i.TgFileID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setImageFileId = `-- name: SetImageFileId :exec
UPDATE post_images
SET tg_file_id = $1
WHERE post = $2
`

type SetImageFileIdParams struct {
	TgFileID string `json:"tg_file_id"`
	Post     string `json:"post"`
}

func (q *Queries) SetImageFileId(ctx context.Context, arg SetImageFileIdParams) error {
	_, err := q.db.Exec(ctx, setImageFileId, arg.TgFileID, arg.Post)
	return err
}
