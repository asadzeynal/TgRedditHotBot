-- name: CreatePostVideo :one
INSERT INTO post_videos (
    post,
    height,
    width,
    duration,
    url,
    audio_url
) VALUES (
$1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetVideosByPost :many
SELECT * FROM post_videos
WHERE post = $1 LIMIT 5;

-- name: SetVideoFileId :exec
UPDATE post_videos
SET tg_file_id = $1
WHERE post = $2;
