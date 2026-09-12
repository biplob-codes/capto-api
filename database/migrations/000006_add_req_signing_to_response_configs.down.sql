ALTER TABLE response_configs
DROP CONSTRAINT check_secret_config_pair,
DROP COLUMN signing_secret,
DROP COLUMN signature_header,
DROP COLUMN tolerance_window;