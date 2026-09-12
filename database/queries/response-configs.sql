-- name: CreateResponseConfig :one
INSERT INTO response_configs
(endpoint_id,method,status_code,headers,body,delay,signing_secret,signature_header,tolerance_window)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING *;

-- name: GetResponseConfigsByEndpointId :many
SELECT * FROM response_configs WHERE endpoint_id=$1;

-- name: UpdateResponseConfig :one
UPDATE response_configs
SET
 method=$1,status_code=$2,headers=$3,body=$4,delay=$5,signing_secret=$6,
 signature_header=$7,tolerance_window=$8
WHERE id=$9
RETURNING *;

-- name: DeleteResponseConfig :exec
DELETE FROM response_configs WHERE id=$1;

-- name: GetResponseConfigForRequest :one
SELECT * FROM response_configs
WHERE endpoint_id=$1 
AND method IN ($2,'ALL')
ORDER BY (method = 'ALL')
LIMIT 1;