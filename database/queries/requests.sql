-- name: CreateRequest :one
WITH update_req_count AS (
  UPDATE endpoints
  SET req_count = req_count + 1
  WHERE endpoints.id = sqlc.arg(endpoint_id)
),
create_req AS (
  INSERT INTO requests (url, remote_addr, body_size, method, duration, headers, body, query_params, endpoint_id)
  VALUES (
    sqlc.arg(url), sqlc.arg(remote_addr), sqlc.arg(body_size), sqlc.arg(method),
    sqlc.arg(duration), sqlc.arg(request_headers), sqlc.arg(request_body),
    sqlc.arg(query_params), sqlc.arg(endpoint_id)
  )
  RETURNING *
),
create_response AS (
  INSERT INTO responses (status_code, headers, body, request_id)
  VALUES (
    sqlc.arg(status_code), sqlc.arg(response_headers), sqlc.arg(response_body),
    (SELECT create_req.id FROM create_req)
  )
  RETURNING *
)
SELECT * FROM create_response;



-- name: AddRequestNote :one
UPDATE requests SET note=$1 WHERE id=$2 RETURNING *;

-- name: GetRequestsByEndpointId :many  
SELECT * FROM requests WHERE endpoint_id=$1;       

-- name: GetRequestById :one
SELECT * FROM requests WHERE id=$1;