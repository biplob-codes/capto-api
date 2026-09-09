CREATE TABLE responses(
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
status_code INT NOT NULL CHECK (status_code BETWEEN 100 AND 599),
headers TEXT,
body TEXT,
request_id UUID,
created_at TIMESTAMPTZ DEFAULT now(),
FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE CASCADE
);