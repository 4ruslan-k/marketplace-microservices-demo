
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS products (
    name varchar NOT NULL,
    id uuid NOT NULL,
    price numeric NOT NULL,
    quantity int NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NULL,
    CONSTRAINT products_pk PRIMARY KEY (id)
);   

