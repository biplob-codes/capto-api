DROP TRIGGER IF EXISTS trg_request_update ON requests;
DROP INDEX IF EXISTS idx_request_endpoint_id;
DROP TABLE requests;
DROP TYPE IF EXISTS req_method;