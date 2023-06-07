
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS products (
    name varchar NOT NULL,
    id uuid NOT NULL,
    CONSTRAINT products_pk PRIMARY KEY (id)
);   

