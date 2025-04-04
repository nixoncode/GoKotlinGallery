-- name: CreateImage :one
INSERT INTO images
    (filename, description, metadata)
VALUES
    ($1, $2, $3)
RETURNING *;

-- name: GetImage :one
SELECT *
FROM images
WHERE filename = $1;

-- name: ListAllImageDetails :many
SELECT filename, metadata, description
FROM images;
