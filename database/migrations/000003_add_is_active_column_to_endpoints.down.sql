ALTER TABLE endpoints 
DROP COLUMN status,
DROP COLUMN req_count;
DROP TYPE IF EXISTS endpoint_status;