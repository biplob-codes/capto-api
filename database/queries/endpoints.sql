-- name: CreateEndpoint :one
INSERT INTO endpoints(label,token)
VALUES ($1,$2) RETURNING *;

-- name: GetEndpoint :one
SELECT * FROM endpoints WHERE id=$1;

-- name: ListEndpoints :many
SELECT * FROM endpoints ORDER BY created_at;

-- name: GetEndpointByToken :one
SELECT * FROM endpoints WHERE token=$1;

-- name: UpdateEndpoint :one
UPDATE endpoints
SET 
 status=$1,
 label=$2
WHERE id=$3
RETURNING *;

-- name: DeleteEndpoint :exec
DELETE FROM endpoints WHERE id=$1;