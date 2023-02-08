-- name: CreatePostVideo :one
INSERT INTO post_videos (
    post,
    height,
    width,
    duration,
    url
) VALUES (
$1, $2, $3, $4, $5
) RETURNING *;

-- name: GetVideosByPost :many
SELECT * FROM post_videos
WHERE post = $1 LIMIT 5;
