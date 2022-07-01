-- name: CreateAuthor :one
INSERT INTO authors (name, bio)
VALUES ($1, $2)
RETURNING *;

-- name: GetAuthor :one
SELECT *
FROM authors
WHERE id = $1
LIMIT 1;

-- name: UpdateAuthor :one
UPDATE authors
SET name = $2,
    bio  = $3
WHERE id = $1
RETURNING *;

-- name: PartialUpdateAuthor :one
UPDATE authors
SET name = CASE WHEN @update_name::boolean THEN @name::VARCHAR(32) ELSE name END,
    bio  = CASE WHEN @update_bio::boolean THEN @bio::TEXT ELSE bio END
WHERE id = @id
RETURNING *;

-- name: DeleteAuthor :exec
DELETE
FROM authors
WHERE id = $1;

-- name: ListAuthors :many
SELECT *
FROM authors
ORDER BY name;

-- name: TruncateAuthor :exec
TRUNCATE authors;