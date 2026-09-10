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
UPDATE requests req
SET note = $1
FROM responses res
WHERE req.id = $2 AND res.request_id = req.id
RETURNING
  req.id,
  req.url,
  req.remote_addr,
  req.body_size,
  req.method,
  req.duration,
  req.note,
  req.headers AS request_headers,
  req.body AS request_body,
  req.query_params,
  req.endpoint_id,
  req.created_at AS request_created_at,
  req.updated_at AS request_updated_at,
  res.id AS response_id,
  res.status_code,
  res.headers AS response_headers,
  res.body AS response_body,
  res.created_at AS response_created_at;

-- name: GetRequestById :one
SELECT
  req.id,
  req.url,
  req.remote_addr,
  req.body_size,
  req.method,
  req.duration,
  req.note,
  req.headers AS request_headers,
  req.body AS request_body,
  req.query_params,
  req.endpoint_id,
  req.created_at AS request_created_at,
  req.updated_at AS request_updated_at,
  res.id AS response_id,
  res.status_code,
  res.headers AS response_headers,
  res.body AS response_body,
  res.created_at AS response_created_at
FROM requests req
JOIN responses res ON res.request_id = req.id
WHERE req.id = $1;

-- name: GetRequestsByEndpointId :many
SELECT
  req.id,
  req.url,
  req.remote_addr,
  req.body_size,
  req.method,
  req.duration,
  req.note,
  req.headers AS request_headers,
  req.body AS request_body,
  req.query_params,
  req.endpoint_id,
  req.created_at AS request_created_at,
  req.updated_at AS request_updated_at,
  res.id AS response_id,
  res.status_code,
  res.headers AS response_headers,
  res.body AS response_body,
  res.created_at AS response_created_at
FROM requests req
JOIN responses res ON res.request_id = req.id
WHERE req.endpoint_id = $1
ORDER BY req.created_at DESC;