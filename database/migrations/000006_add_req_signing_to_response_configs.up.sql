ALTER TABLE response_configs
ADD COLUMN signing_secret TEXT,
ADD COLUMN signature_header TEXT,
ADD COLUMN tolerance_window INT NOT NULL DEFAULT 300,
ADD CONSTRAINT check_secret_config_pair CHECK (
  (signing_secret IS NOT NULL AND signature_header IS NOT NULL) OR
  (signing_secret IS NULL AND signature_header IS NULL)
);