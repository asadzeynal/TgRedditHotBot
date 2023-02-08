-- name: CreatePostImage :one
INSERT INTO post_images (
    post,
    url
) VALUES (
$1, $2
) RETURNING *;

-- name: GetImagesByPost :many
SELECT * FROM post_images
WHERE post = $1 LIMIT 10;
