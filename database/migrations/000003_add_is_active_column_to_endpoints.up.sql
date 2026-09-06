CREATE TYPE endpoint_status AS ENUM ('ACTIVE','INACTIVE');

ALTER TABLE endpoints 
ADD COLUMN status endpoint_status DEFAULT 'ACTIVE',
ADD COLUMN req_count INT DEFAULT 0;