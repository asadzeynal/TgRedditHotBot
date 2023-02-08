-- name: CreatePost :one
INSERT INTO posts (
    id,
    title,
    url
) VALUES (
$1, $2, $3
) RETURNING *;

-- name: GetRandomPost :one
SELECT * FROM posts
TABLESAMPLE system_rows(1);