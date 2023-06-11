
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS customers (
    id uuid NOT NULL,
    name varchar NOT NULL,
    email varchar NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NULL,
    CONSTRAINT customers_pk PRIMARY KEY (id),
    CONSTRAINT customers_un UNIQUE (email)
);   

