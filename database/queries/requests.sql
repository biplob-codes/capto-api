-- name: CreateRequest :one
WITH update_req_count AS (
UPDATE endpoints 
SET req_count=req_count+1 
WHERE id=$9
 )
INSERT INTO requests(url,remote_addr,body_size,method,duration,headers,body,query_params,endpoint_id)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING *;

-- name: AddRequestNote :one
UPDATE requests SET note=$1 WHERE id=$2 RETURNING *;

-- name: GetRequestsByEndpointId :many  
SELECT * FROM requests WHERE endpoint_id=$1;       

-- name: GetRequestById :one
SELECT * FROM requests WHERE id=$1;