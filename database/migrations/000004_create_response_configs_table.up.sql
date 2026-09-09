CREATE TYPE res_method AS ENUM ('ALL','GET','POST','PUT','PATCH','DELETE','OPTIONS');

CREATE TABLE response_configs(
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
method res_method NOT NULL DEFAULT 'ALL',
endpoint_id UUID NOT NULL,
delay INT DEFAULT 0,
headers TEXT DEFAULT '{"Content-Type": "application/json"}',
body TEXT DEFAULT '{"received": true}',
status_code INT DEFAULT 200 CHECK (status_code BETWEEN 100 AND 599),
created_at TIMESTAMPTZ DEFAULT now(),
updated_at TIMESTAMPTZ DEFAULT now(),
UNIQUE (endpoint_id,method),
FOREIGN KEY (endpoint_id) REFERENCES endpoints(id) ON DELETE CASCADE
);

CREATE TRIGGER trg_response_configs_update 
BEFORE UPDATE ON response_configs
FOR EACH ROW EXECUTE FUNCTION set_updated_at();