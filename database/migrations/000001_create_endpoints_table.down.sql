DROP TRIGGER IF EXISTS trg_endpoint_update ON endpoints;
DROP FUNCTION IF EXISTS set_updated_at();
DROP INDEX IF EXISTS idx_endpoint_token;
DROP TABLE IF EXISTS endpoints;