CREATE TYPE req_method AS ENUM ('GET','POST','PUT','PATCH','DELETE','OPTIONS');

CREATE TABLE requests(
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
url TEXT NOT NULL,
remote_addr TEXT NOT NULL,
body_size INT NOT NULL,
method req_method NOT NULL,
duration INT NOT NULL,
note TEXT,
headers TEXT,
body TEXT,
query_params TEXT,
endpoint_id UUID NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
FOREIGN KEY (endpoint_id) REFERENCES endpoints(id) ON DELETE CASCADE
);

CREATE INDEX idx_request_endpoint_id ON requests(endpoint_id);


CREATE TRIGGER trg_request_update
BEFORE UPDATE ON requests
FOR EACH ROW EXECUTE FUNCTION set_updated_at();