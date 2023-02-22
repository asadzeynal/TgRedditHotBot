-- name: CreatePostImage :one
INSERT INTO post_images (
    post,
    url,
    is_gif
) VALUES (
$1, $2, $3
) RETURNING *;

-- name: GetImagesByPost :many
SELECT * FROM post_images
WHERE post = $1 LIMIT 10;
