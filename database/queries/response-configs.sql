-- name: CreateResponseConfig :one
INSERT INTO response_configs
(endpoint_id,method,status_code,headers,body,delay)
VALUES ($1,$2,$3,$4,$5,$6) RETURNING *;