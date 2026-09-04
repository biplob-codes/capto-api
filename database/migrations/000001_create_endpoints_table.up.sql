CREATE TABLE endpoints(
id UUID PRIMARY KEY DEFAULT gen_random_uuid() ,
label TEXT NOT NULL,
token TEXT NOT NULL UNIQUE,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT  now()
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
NEW.updated_at=now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_endpoint_update
BEFORE UPDATE ON endpoints
FOR EACH ROW EXECUTE FUNCTION set_updated_at();