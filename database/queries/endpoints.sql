-- name: CreateEndpoint :one
INSERT INTO endpoints(label,token)
VALUES ($1,$2) RETURNING *;

-- name: GetEndpoint :one
SELECT * FROM endpoints WHERE id=$1;

-- name: ListEndpoints :many
SELECT * FROM endpoints ORDER BY created_at;

