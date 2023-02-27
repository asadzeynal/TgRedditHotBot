// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: post_videos.sql

package db

import (
	"context"
)

const createPostVideo = `-- name: CreatePostVideo :one
INSERT INTO post_videos (
    post,
    height,
    width,
    duration,
    url
) VALUES (
$1, $2, $3, $4, $5
) RETURNING id, post, height, width, duration, url, tg_file_id
`

type CreatePostVideoParams struct {
	Post     string `json:"post"`
	Height   int32  `json:"height"`
	Width    int32  `json:"width"`
	Duration int32  `json:"duration"`
	Url      string `json:"url"`
}

func (q *Queries) CreatePostVideo(ctx context.Context, arg CreatePostVideoParams) (PostVideo, error) {
	row := q.db.QueryRow(ctx, createPostVideo,
		arg.Post,
		arg.Height,
		arg.Width,
		arg.Duration,
		arg.Url,
	)
	var i PostVideo
	err := row.Scan(
		&i.ID,
		&i.Post,
		&i.Height,
		&i.Width,
		&i.Duration,
		&i.Url,
		&i.TgFileID,
	)
	return i, err
}

const getVideosByPost = `-- name: GetVideosByPost :many
SELECT id, post, height, width, duration, url, tg_file_id FROM post_videos
WHERE post = $1 LIMIT 5
`

func (q *Queries) GetVideosByPost(ctx context.Context, post string) ([]PostVideo, error) {
	rows, err := q.db.Query(ctx, getVideosByPost, post)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PostVideo{}
	for rows.Next() {
		var i PostVideo
		if err := rows.Scan(
			&i.ID,
			&i.Post,
			&i.Height,
			&i.Width,
			&i.Duration,
			&i.Url,
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
